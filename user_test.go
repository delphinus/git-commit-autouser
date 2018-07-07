package main

import (
	"errors"
	"os/user"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetFromConfig(t *testing.T) {
	a := assert.New(t)
	u, err := user.Current()
	a.NoError(err)
	for _, c := range []struct {
		name      string
		key, val  string
		want      User
		errPrefix string
	}{
		{
			name:      "invalid url-regexp",
			key:       `url-regexp`,
			val:       `hoge(hoge`,
			errPrefix: "error in regexp in ",
		},
		{
			name: "valid url-regexp",
			key:  `url-regexp`,
			val:  `github\.example\.com`,
			want: User{URLRegexp: regexp.MustCompile(`github\.example\.com`)},
		},
		{
			name: "valid name",
			key:  `name`,
			val:  `Foo Bar`,
			want: User{Name: []byte(`Foo Bar`)},
		},
		{
			name: "valid email",
			key:  `email`,
			val:  `foo@example.com`,
			want: User{Email: []byte(`foo@example.com`)},
		},
		{
			name: "valid email",
			key:  `email`,
			val:  `foo@example.com`,
			want: User{Email: []byte(`foo@example.com`)},
		},
		{
			name: "valid hub-config",
			key:  `hub-config`,
			val:  `~/.config/hub`,
			want: User{HubConfig: []byte(u.HomeDir + `/.config/hub`)},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			u := User{}
			err := u.setFromConfig([]byte(c.key), []byte(c.val))
			if c.errPrefix == "" {
				a.NoError(err)
				a.Equal(c.want, u)
			} else {
				a.Error(err)
				a.True(strings.HasPrefix(err.Error(), c.errPrefix))
			}
		})
	}
}

func TestConfigUsers(t *testing.T) {
	for _, c := range []struct {
		name      string
		outBytes  []byte
		outErr    error
		errPrefix string
		users     Users
	}{
		{
			name:      "if command output none, show instruction",
			outErr:    errors.New("dummy error"),
			errPrefix: "show instruction",
		},
		{
			name:      "if command error, return error",
			outBytes:  []byte(`dummy`),
			outErr:    errors.New("hoge"),
			errPrefix: "hoge",
		},
		{
			name:      "if no error but invalid config, return error",
			outBytes:  []byte(`hogehoge`),
			errPrefix: "invalid config from git",
		},
		{
			name:      "if no error but no line, show instruction",
			errPrefix: "show instruction",
		},
		{
			name: "if no error and valid output, return valid users",
			outBytes: []byte(`autouser-ghe.url-regexp git\.example\.com
autouser-ghe.name Foo Bar
autouser-gitlab.url-regexp gitlab\.com
autouser-gitlab.name Bar Foo
`),
			users: Users{
				"ghe": {
					Name:      []byte("Foo Bar"),
					URLRegexp: regexp.MustCompile(`git\.example\.com`),
				},
				"gitlab": {
					Name:      []byte("Bar Foo"),
					URLRegexp: regexp.MustCompile(`gitlab\.com`),
				},
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			a := assert.New(t)
			defer ReplaceExecCommand(c.outBytes, c.outErr)()
			users, err := configUsers()
			if c.errPrefix != "" {
				a.Error(err)
				a.True(strings.HasPrefix(err.Error(), c.errPrefix))
				return
			}
			a.NoError(err)
			a.Equal(c.users, users)
		})
	}
}
