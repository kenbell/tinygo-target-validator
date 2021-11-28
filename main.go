package main

import (
	"fmt"
	"os"
	"path"
)

const (
	appName        = "tinygo-target-validator"
	dirPeripherals = "./tests/peripherals"
	dirFeatures    = "./tests/features"
)

type featureTestResult struct {
	Target  string `json:"target"`
	Feature string `json:"feature"`
	Passed  bool   `json:"passed"`
	Output  string `json:"output"`
}

type peripheralTestResult struct {
	Target          string `json:"target"`
	PeripheralClass string `json:"pclass"`
	Feature         string `json:"feature"`
	Peripheral      string `json:"peripheral"`
	Passed          bool   `json:"passed"`
	Output          string `json:"output"`
}

type testReport struct {
	Features    []featureTestResult    `json:"features"`
	Peripherals []peripheralTestResult `json:"peripherals"`
}

type testParams struct {
	Peripheral string
}

func main() {
	testCmd := newTestCommand()

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {

	case "test":
		testCmd.Parse(os.Args[2:])
		testCmd.Execute()
	case "help":
		usage()
		os.Exit(0)
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Printf("Syntax for %s:\n", path.Base(os.Args[0]))
	fmt.Printf("  %s <command> [flags...]\n", path.Base(os.Args[0]))
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  test    performs the tests to validate tinygo targets")
	fmt.Println("  help    shows this help")
	fmt.Println()
	fmt.Printf("To get help on each command do: %s <command> -h\n", os.Args[0])
}
