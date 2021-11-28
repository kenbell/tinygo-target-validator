package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

// This file contains the code to build test apps to check for
// compilation errors.

func buildFromTemplate(test *template.Template, target string, peripheral string) (string, error) {
	tmpDir, err := ioutil.TempDir("", appName+"-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	sourcePath := path.Join(tmpDir, "main.go")
	source, err := os.Create(sourcePath)
	if err != nil {
		return "", fmt.Errorf("unable to output temp Go source %s: %w", sourcePath, err)
	}

	err = test.Execute(source, testParams{Peripheral: peripheral})
	source.Close()
	if err != nil {
		return "", fmt.Errorf("unable to execute Go template: %w", err)
	}

	cmd := exec.Command("tinygo", "build", "-target="+target, "-opt=0", "-o", "test.elf", sourcePath)

	output := bytes.Buffer{}
	cmd.Stdout = &output
	cmd.Stderr = &output
	err = cmd.Start()
	if err != nil {
		return "", fmt.Errorf("unable to start '%s', %w", cmd.String(), err)
	}

	err = cmd.Wait()
	if err != nil {
		return output.String(), fmt.Errorf("tinygo build '%s' failed, %w", cmd.String(), err)
	}

	return output.String(), nil
}
