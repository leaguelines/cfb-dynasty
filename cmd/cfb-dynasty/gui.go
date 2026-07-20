package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func runGUI(args []string) int {
	path, err := exec.LookPath("cfb-dynasty-gui")
	if err != nil {
	fmt.Fprint(os.Stderr, `gui: cfb-dynasty-gui not found on PATH.

Build the desktop app (from repo root):
  go build -tags production -o cfb-dynasty-gui ./cmd/cfb-dynasty-gui
  go run -tags production ./cmd/cfb-dynasty-gui --schema-dir ./data/schemas

  # or package with the Wails CLI:
  cd cmd/cfb-dynasty-gui && wails build

First launch: choose your schema folder (C27_*.gz), then open a dynasty save.
On Windows, saves under Documents\EA SPORTS CFB27\saves are listed automatically.
`)
	return 2
}

	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
			return 1
		}
		fmt.Fprintf(os.Stderr, "gui: %v\n", err)
		return 1
	}
	return 0
}
