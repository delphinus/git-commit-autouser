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
