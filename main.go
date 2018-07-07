package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var (
	nocolor bool
)

func main() {
	flag.BoolVar(&nocolor, "nocolor", false, "show output without color")
	flag.Parse()
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
	var args []string
	if !nocolor {
		args = []string{"-c", "color.status=always"}
	}
	args = append(args, "commit")
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	return cmd.Run()
}
