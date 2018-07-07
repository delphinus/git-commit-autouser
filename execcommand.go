package main

import "os/exec"

var (
	execCommand func(string, ...string) outputter
)

func init() {
	execCommand = func(name string, args ...string) outputter {
		return exec.Command(name, args...)
	}
}
