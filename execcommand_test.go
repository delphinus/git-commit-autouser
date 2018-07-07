package main

type testOutputter struct {
	out []byte
	err error
}

func (to *testOutputter) Output() ([]byte, error) {
	return to.out, to.err
}

// ReplaceExecCommand replaces exec.Command to test func
func ReplaceExecCommand(out []byte, err error) func() {
	original := execCommand
	execCommand = func(string, ...string) outputter {
		return &testOutputter{out, err}
	}
	return func() { execCommand = original }
}
