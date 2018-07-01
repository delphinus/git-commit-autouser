package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

const (
	gitConfigPrefix = "autouser-"
	gitConfigRegexp = gitConfigPrefix + `.+\.`
)

var (
	origin = []byte("origin")
	prefix = []byte(gitConfigPrefix)
	tilde  = []byte{'~'}
	home   = []byte(os.Getenv("HOME"))
)

// User has config for several domains
type User struct {
	Name, Email, HubConfig []byte
	URLRegexp              *regexp.Regexp
}

func (u *User) setFromConfig(line []byte) error {
	kv := bytes.SplitN(line, []byte{' '}, 2)
	if len(kv) != 2 {
		return fmt.Errorf("invalid config from git: %s", line)
	}
	fullKeyName := bytes.TrimPrefix(kv[0], prefix)
	val := kv[1]
	nameKey := bytes.SplitN(fullKeyName, []byte{'.'}, 2)
	switch key := nameKey[1]; {
	case bytes.Equal(key, []byte("url-regexp")):
		re, err := regexp.Compile(string(key))
		if err != nil {
			return fmt.Errorf("error in regexp in %s: %v", key, err)
		}
		u.URLRegexp = re
	case bytes.Equal(key, []byte("name")):
		u.Name = val
	case bytes.Equal(key, []byte("email")):
		u.Email = val
	case bytes.Equal(key, []byte("hub-config")):
		u.HubConfig = replaceTilde(val)
	}
	return nil
}

// Env returns env for the user
func (u *User) Env() []string {
	return []string{
		"GIT_COMMITTER_NAME=" + string(u.Name),
		"GIT_COMMITTER_EMAIL=" + string(u.Email),
		"GIT_AUTHOR_NAME=" + string(u.Name),
		"GIT_AUTHOR_EMAIL=" + string(u.Email),
		"HUB_CONFIG=" + string(u.HubConfig),
	}
}

func configUsers() (Users, error) {
	cmd := exec.Command("git", "config", "--get-regexp", gitConfigRegexp)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var users Users
	for _, line := range bytes.Split(out, []byte{'\n'}) {
		u := User{}
		if err := u.setFromConfig(line); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if len(users) == 0 {
		return nil, ErrShowInstruction{}
	}
	return users, nil
}

func replaceTilde(path []byte) []byte {
	if !bytes.HasPrefix(path, tilde) {
		return path
	}
	return bytes.Replace(path, tilde, home, 1)
}
