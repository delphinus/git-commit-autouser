package main

import (
	"bytes"
	"os/exec"
)

var (
	pushURLPrefix = []byte("  Push  URL: ")
)

// Users is a bunch of User
type Users []User

// Env returns env setting for users
func (us Users) Env() ([]string, error) {
	url, err := originRemoteURL()
	if err != nil {
		return nil, err
	}
	for _, u := range us {
		if u.URLRegexp.Match(url) {
			return u.Env(), nil
		}
	}
	return nil, ErrNotMatch{url}
}

func originRemoteURL() ([]byte, error) {
	cmd := exec.Command("git", "remote", "show", "-n", string(origin))
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var url []byte
	for _, line := range bytes.Split(out, []byte{'\n'}) {
		if bytes.Contains(line, pushURLPrefix) {
			url = bytes.TrimPrefix(line, pushURLPrefix)
			break
		}
	}
	if bytes.Equal(url, origin) {
		return nil, ErrNotSetOrigin{}
	}
	return url, nil
}
