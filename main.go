package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var (
	flagNocolor, flagSelfupdate, flagVersion bool
)

func main() {
	flag.BoolVar(&flagNocolor, "nocolor", false, "show output without color")
	flag.BoolVar(&flagSelfupdate, "selfupdate", false, "self-update the binary")
	flag.BoolVar(&flagVersion, "version", false, "show the version")
	flag.BoolVar(&flagVersion, "v", false, "show the version (shorthand)")
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
	if flagSelfupdate {
		if err := doSelfupdate(); err != nil {
			return err
		}
		return nil
	}
	if flagVersion {
		showVersion()
		return nil
	}
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
	if !flagNocolor {
		args = []string{"-c", "color.status=always"}
	}
	args = append(args, "commit")
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	envs := append(os.Environ(), env...)
	cmd.Env = envs
	return cmd.Run()
}
