package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

var (
	flagNocolor, flagSelfupdate, flagVersion bool

	osArgs = os.Args
	flags  = map[string]bool{
		"nocolor":    true,
		"selfupdate": true,
		"version":    true,
		"v":          true,
	}
)

func main() {
	// Suppress any output to avoid help message duplication.
	flag.CommandLine.Init(osArgs[0], flag.ContinueOnError)
	flag.CommandLine.SetOutput(ioutil.Discard)
	flag.BoolVar(&flagNocolor, "nocolor", false, "show output without color")
	flag.BoolVar(&flagSelfupdate, "selfupdate", false, "self-update the binary")
	flag.BoolVar(&flagVersion, "version", false, "show the version")
	flag.BoolVar(&flagVersion, "v", false, "show the version (shorthand)")
	if err := flag.CommandLine.Parse(osArgs[1:]); err == flag.ErrHelp {
		flag.CommandLine.SetOutput(os.Stderr)
		flag.Usage()
		os.Exit(1)
	}
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
		return doSelfupdate()
	}
	if flagVersion {
		return showVersion()
	}
	users, err := configUsers()
	if err != nil {
		return err
	}
	env, info, err := users.Env()
	if err != nil {
		w, ok := err.(WarningMessager)
		if !ok {
			return err
		}
		fmt.Fprintf(os.Stderr, "[warning]: %s\n", w.WarningMessage())
	}
	fmt.Printf("user: %s\n", info)
	return run(env)
}

func run(env []string) error {
	var args []string
	if !flagNocolor {
		args = []string{"-c", "color.status=always"}
	}
	args = append(args, "commit")
	args = append(args, argsWithoutKnownFlags()...)
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	envs := append(os.Environ(), env...)
	cmd.Env = envs
	return cmd.Run()
}

func argsWithoutKnownFlags() (args []string) {
	for _, a := range osArgs[1:] {
		if a[0] == '-' {
			var f string
			if a[1] == '-' {
				f = a[2:]
			} else {
				f = a[1:]
			}
			if _, ok := flags[f]; ok {
				continue
			}
		}
		args = append(args, a)
	}
	return
}
