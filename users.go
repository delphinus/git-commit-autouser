package main

import (
	"bytes"
	"errors"
	"fmt"
)

type outputter interface {
	Output() ([]byte, error)
}

var (
	pushURLPrefix = []byte("  Push  URL: ")
)

// Users is a bunch of User
type Users map[string]*User

// Env returns env setting for users
func (us Users) Env() ([]string, string, error) {
	url, err := originRemoteURL()
	if err != nil {
		return nil, "", err
	}
	for _, u := range us {
		if u.URLRegexp.Match(url) {
			return u.Env(), fmt.Sprintf("%s <%s>", u.Name, u.Email), nil
		}
	}
	return nil, "", ErrNotMatch{url}
}

// User returns the user with supplied name
func (us Users) User(name []byte) *User {
	key := string(name)
	if u, ok := us[key]; ok {
		return u
	}
	us[key] = &User{}
	return us[key]
}

func originRemoteURL() ([]byte, error) {
	cmd := execCommand("git", "remote", "show", "-n", string(origin))
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
	if url == nil {
		return nil, errors.New("invalid output from git remote")
	} else if bytes.Equal(url, origin) {
		return nil, ErrNotSetOrigin{}
	}
	return url, nil
}
