package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	outBytes []byte
	outErr   error
)

type testOutputter struct{}

func (testOutputter) Output() ([]byte, error) {
	return outBytes, outErr
}

func TestOriginRemoteURL(t *testing.T) {
	execCommand = func(_ string, _ ...string) outputter {
		return testOutputter{}
	}
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
			outBytes = c.outBytes
			outErr = c.outErr
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
