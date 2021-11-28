package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// This file contains the code to compile a list of valid targets to
// test.

func detectTargets() ([]string, error) {
	cmd := exec.Command("tinygo", "targets")
	output := bytes.Buffer{}
	cmd.Stdout = &output
	cmd.Stderr = &output
	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("unable to start '%s', %w", cmd.String(), err)
	}

	err = cmd.Wait()
	if err != nil {
		return nil, fmt.Errorf("targets command '%s' failed, %w", cmd.String(), err)
	}

	listStr := strings.TrimSpace(output.String())
	targets := strings.Split(strings.ReplaceAll(listStr, "\r\n", "\n"), "\n")

	return targets, nil
}
