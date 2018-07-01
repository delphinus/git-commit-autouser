package main

import "fmt"

// ErrorMessager can return message to show
type ErrorMessager interface {
	ErrorMessage() string
}

// WarningMessager can return warning message
type WarningMessager interface {
	WarningMessage() string
}

// ErrShowInstruction indicates caller should show the instruction.
type ErrShowInstruction struct{}

func (e ErrShowInstruction) Error() string { return "show instruction" }

// ErrorMessage satisfies Messager
func (e ErrShowInstruction) ErrorMessage() string {
	return fmt.Sprintf(`
No user setting found. You should add to ~/.gitconfig like the following.
-----
[%sgithub]
  url-regexp = github\\.com
  name = "Foo Bar"
  email = foo@private.com
[%sghe]
  url-regexp = git\\.company\\.com
  name = "Foo Bar"
  email = bar@company.com
-----

`, gitConfigPrefix, gitConfigPrefix)
}

// ErrNotSetOrigin indicates user should set remote URL for origin
type ErrNotSetOrigin struct{}

func (e ErrNotSetOrigin) Error() string { return "not set origin" }

// WarningMessage satisfies WarningMessager
func (e ErrNotSetOrigin) WarningMessage() string {
	return "remote `origin` is not configured"
}

// ErrNotMatch indicates .gitconfig has no setting matching the URL
type ErrNotMatch struct{ URL []byte }

func (e ErrNotMatch) Error() string { return "not match" }

// WarningMessage satisfies WarningMessager
func (e ErrNotMatch) WarningMessage() string {
	return "no setting for this remote URL: " + string(e.URL)
}
