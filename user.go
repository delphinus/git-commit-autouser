package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/user"
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
)

// User has config for several domains
type User struct {
	Name, Email, HubConfig []byte
	URLRegexp              *regexp.Regexp
}

func (u *User) setFromConfig(key, val []byte) (err error) {
	switch {
	case bytes.Equal(key, []byte("url-regexp")):
		u.URLRegexp, err = regexp.Compile(string(val))
		if err != nil {
			return fmt.Errorf("error in regexp in %s: %v", key, err)
		}
	case bytes.Equal(key, []byte("name")):
		u.Name = val
	case bytes.Equal(key, []byte("email")):
		u.Email = val
	case bytes.Equal(key, []byte("hub-config")):
		u.HubConfig, err = replaceTilde(val)
		if err != nil {
			return fmt.Errorf("error in replacing tilde: %v", err)
		}
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
	cmd := execCommand("git", "config", "--get-regexp", gitConfigRegexp)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	users := Users{}
	for _, line := range bytes.Split(out, []byte{'\n'}) {
		kv := bytes.SplitN(line, []byte{' '}, 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid config from git: %s", line)
		}
		fullKeyName := bytes.TrimPrefix(kv[0], prefix)
		val := kv[1]
		nameKey := bytes.SplitN(fullKeyName, []byte{'.'}, 2)
		u := users.User(nameKey[0])
		if err := u.setFromConfig(nameKey[1], val); err != nil {
			return nil, err
		}
	}
	if len(users) == 0 {
		return nil, ErrShowInstruction{}
	}
	return users, nil
}

func replaceTilde(path []byte) ([]byte, error) {
	if !bytes.HasPrefix(path, tilde) {
		return path, nil
	}
	u, err := user.Current()
	if err != nil {
		return nil, errors.New("cannot get the current user")
	}
	return bytes.Replace(path, tilde, []byte(u.HomeDir), 1), nil
}
