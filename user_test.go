package main

import (
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
		line      string
		want      User
		errPrefix string
	}{
		{
			name:      "invalid config",
			line:      `hoge`,
			errPrefix: "invalid config from git:",
		},
		{
			name:      "invalid url-regexp",
			line:      `autouser-ghe.url-regexp hoge(hoge`,
			errPrefix: "error in regexp in ",
		},
		{
			name: "valid url-regexp",
			line: `autouser-ghe.url-regexp github\.example\.com`,
			want: User{URLRegexp: regexp.MustCompile(`github\.example\.com`)},
		},
		{
			name: "valid name",
			line: `autouser-ghe.name Foo Bar`,
			want: User{Name: []byte(`Foo Bar`)},
		},
		{
			name: "valid email",
			line: `autouser-ghe.email foo@example.com`,
			want: User{Email: []byte(`foo@example.com`)},
		},
		{
			name: "valid email",
			line: `autouser-ghe.email foo@example.com`,
			want: User{Email: []byte(`foo@example.com`)},
		},
		{
			name: "valid hub-config",
			line: `autouser-ghe.hub-config ~/.config/hub`,
			want: User{HubConfig: []byte(u.HomeDir + `/.config/hub`)},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			u := User{}
			err := u.setFromConfig([]byte(c.line))
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
