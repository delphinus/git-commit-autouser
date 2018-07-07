package main

import (
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOriginRemoteURL(t *testing.T) {
	a := assert.New(t)
	for _, c := range []struct {
		name        string
		outBytes    []byte
		outErr      error
		expected    string
		expectedErr string
	}{
		{
			name:        "if Output returns error, returns error",
			outErr:      errors.New("hoge"),
			expectedErr: "hoge",
		},
		{
			name:        "if invalid output, return error",
			outBytes:    []byte(`hogehogeo`),
			expectedErr: "invalid output from git remote",
		},
		{
			name: "if URL is origin, return ErrNotSetOrigin",
			outBytes: []byte(`* remote origin
  Fetch URL: origin
  Push  URL: origin
  HEAD branch: (not queried)
  Local ref configured for 'git push' (status not queried):
    (matching) pushes to (matching)
`),
			expectedErr: ErrNotSetOrigin{}.Error(),
		},
		{
			name: "if URL is valid, return URL",
			outBytes: []byte(`* remote origin
  Fetch URL: https://ghe.example.com/hoge/fuga
  Push  URL: ssh://git@ghe.example.com/hoge/fuga
  HEAD branch: (not queried)
  Remote branches: (status not queried)
    development
    master
  Local branches configured for 'git pull':
    development merges with remote development
    master      merges with remote master
  Local ref configured for 'git push' (status not queried):
    (matching) pushes to (matching)
`),
			expected: "ssh://git@ghe.example.com/hoge/fuga",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			defer ReplaceExecCommand(c.outBytes, c.outErr)()
			url, err := originRemoteURL()
			if c.expectedErr != "" {
				a.EqualError(err, c.expectedErr)
				return
			}
			a.NoError(err)
			a.Equal(c.expected, string(url))
		})
	}
}

func TestEnv(t *testing.T) {
	a := assert.New(t)
	re, err := regexp.Compile(`ghe\.example\.com`)
	a.NoError(err)
	us := Users{
		{
			Name:      []byte(`Foo Bar`),
			Email:     []byte(`foo@example.com`),
			URLRegexp: re,
		},
	}
	for _, c := range []struct {
		name        string
		outBytes    []byte
		outErr      error
		expected    []string
		expectedErr string
	}{
		{
			name:        "if originRemoteURL returns error, returns error",
			outErr:      errors.New("hoge"),
			expectedErr: "hoge",
		},
		{
			name: "if url not match, returns error",
			outBytes: []byte(`* remote origin
  Fetch URL: https://hoge.hogeo.com
  Push  URL: https://hoge.hogeo.com
  HEAD branch: (not queried)
  Local ref configured for 'git push' (status not queried):
    (matching) pushes to (matching)
`),
			expectedErr: ErrNotMatch{[]byte(`hoge.hogeo.com`)}.Error(),
		},
		{
			name: "if url match, returns valid env",
			outBytes: []byte(`* remote origin
  Fetch URL: https://ghe.example.com/hoge/fuga
  Push  URL: ssh://git@ghe.example.com/hoge/fuga
  HEAD branch: (not queried)
  Remote branches: (status not queried)
    development
    master
  Local branches configured for 'git pull':
    development merges with remote development
    master      merges with remote master
  Local ref configured for 'git push' (status not queried):
    (matching) pushes to (matching)
`),
			expected: []string{
				"GIT_COMMITTER_NAME=Foo Bar",
				"GIT_COMMITTER_EMAIL=foo@example.com",
				"GIT_AUTHOR_NAME=Foo Bar",
				"GIT_AUTHOR_EMAIL=foo@example.com",
				"HUB_CONFIG=",
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			defer ReplaceExecCommand(c.outBytes, c.outErr)()
			env, err := us.Env()
			if c.expectedErr != "" {
				a.EqualError(err, c.expectedErr)
				return
			}
			a.Equal(c.expected, env)
		})
	}
}
