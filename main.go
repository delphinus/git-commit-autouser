package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if err := process(); err != nil {
		if m, ok := err.(ErrorMessager); ok {
			fmt.Fprintf(os.Stderr, m.ErrorMessage())
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func process() error {
	users, err := configUsers()
	if err != nil {
		return err
	}
	env, err := users.Env()
	if err != nil {
		w, ok := err.(WarningMessager)
		if !ok {
			return err
		}
		fmt.Fprintf(os.Stderr, "[warning]: %s\n", w.WarningMessage())
	}
	return run(env)
}

func run(env []string) error {
	cmd := exec.Command("git", "commit")
	cmd.Env = env
	return cmd.Run()
}
