package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// This file contains the code to execute tests against targets.

type testCmd struct {
	*flag.FlagSet
	target *string
	pclass *string
}

type testParameters struct {
	Targets           []string
	PeripheralClasses []string
}

func newTestCommand() *testCmd {
	cmd := testCmd{}
	cmd.FlagSet = flag.NewFlagSet("test", flag.ExitOnError)
	cmd.target = cmd.FlagSet.String("target", "", "limits tests to specified target")
	cmd.pclass = cmd.FlagSet.String("pclass", "", "limits tests to specified peripheral class")
	return &cmd
}

func (cmd *testCmd) Execute() error {
	report := testReport{}

	// Determine the list of targets
	var targets []string
	var err error
	if *cmd.target != "" {
		targets = strings.Split(*cmd.target, ",")
	} else {
		targets, err = detectTargets()
		if err != nil {
			return fmt.Errorf("failed to get list of targets: %w", err)
		}
	}
	fmt.Printf("targets: %s\n", strings.Join(targets, ","))

	// Run the feature tests (unless only tests for specific peripheral requested)
	if *cmd.pclass == "" {
		featureResults, err := cmd.featureTests(targets)
		if err != nil {
			return err
		}
		report.Features = append(report.Features, featureResults...)
	}

	// Run the peripheral tests
	peripheralResults, err := cmd.peripheralTests(targets)
	if err != nil {
		return err
	}
	report.Peripherals = append(report.Peripherals, peripheralResults...)

	// Output the report
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	fmt.Printf("results: %s\n", string(data))

	return nil
}

func (cmd *testCmd) peripheralTests(targets []string) ([]peripheralTestResult, error) {
	allResults := make([]peripheralTestResult, 0, 1000)

	var pClasses []string
	if *cmd.pclass != "" {
		pClasses = strings.Split(*cmd.pclass, ",")
	} else {
		pClassesDirs, err := ioutil.ReadDir(dirPeripherals)
		if err != nil {
			return nil, fmt.Errorf("failed to enumerate peripherals directory looking for tests: %w", err)
		}

		pClasses = make([]string, 0, len(pClassesDirs))
		for _, pClassDir := range pClassesDirs {
			pClasses = append(pClasses, pClassDir.Name())
		}
	}
	fmt.Printf("peripheral classes: %s\n", strings.Join(pClasses, ","))

	for _, pClass := range pClasses {
		classTests, err := loadPeripheralTests(pClass)
		if err != nil {
			return nil, fmt.Errorf("failed to load tests for %s: %w", pClass, err)
		}

		for _, target := range targets {
			// Detect this targets peripherals of this class once, use for multiple tests
			fmt.Printf("detecting %s %s peripherals ... ", target, pClass)
			peripherals, err := detectPeripherals(target, pClass)
			if err != nil {
				return nil, fmt.Errorf("unable to detect peripherals of class %s for target %s: %w", pClass, target, err)
			}
			fmt.Println(strings.Join(peripherals, " "))

			// Run all tests for this peripheral class
			fmt.Printf("testing %s %s peripherals ...", target, pClass)
			for name, test := range classTests {
				results, err := runPeripheralTestForTarget(pClass, name, test, target, peripherals)
				if err != nil {
					return nil, fmt.Errorf("unable to run test '%s.%s' for target %s: %w", pClass, test.Name(), target, err)
				}

				fmt.Printf(" %s", name)

				allResults = append(allResults, results...)
			}

			fmt.Println()
		}
	}

	return allResults, nil
}

func (cmd *testCmd) featureTests(targets []string) ([]featureTestResult, error) {
	allResults := make([]featureTestResult, 0, 1000)

	tests, err := loadFeatureTests()
	if err != nil {
		return nil, fmt.Errorf("failed to load feature tests: %w", err)
	}

	for _, target := range targets {
		fmt.Printf("testing %s features ...", target)

		for name, test := range tests {
			result, err := runTestForFeature(name, test, target)
			if err != nil {
				return nil, fmt.Errorf("unable to run test '%s' for target %s: %w", test.Name(), target, err)
			}

			fmt.Printf(" %s", name)

			allResults = append(allResults, result)
		}

		fmt.Println()
	}

	return allResults, nil
}

func loadFeatureTests() (map[string]*template.Template, error) {
	featureTests := map[string]*template.Template{}
	tests, err := ioutil.ReadDir(dirFeatures)
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate feature tests: %w", err)
	}

	for _, test := range tests {
		tmpl, err := loadFeatureTest(test.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to load test %s: %w", test.Name(), err)
		}

		featureTests[test.Name()] = tmpl
	}

	return featureTests, nil
}

func loadPeripheralTests(pClass string) (map[string]*template.Template, error) {
	classTests := map[string]*template.Template{}
	tests, err := ioutil.ReadDir(path.Join(dirPeripherals, pClass))
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate tests for '%s' peripheral: %w", pClass, err)
	}

	for _, test := range tests {
		tmpl, err := loadPeripheralTest(pClass, test.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to load test %s.%s: %w", pClass, test.Name(), err)
		}

		classTests[test.Name()] = tmpl
	}

	return classTests, nil
}

func loadFeatureTest(test string) (*template.Template, error) {
	testPath := path.Join(dirFeatures, test)
	_, err := os.Stat(testPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("test '%s' not found", test)
	} else if err != nil {
		return nil, fmt.Errorf("unable to open '%s' directory as test: %w", testPath, err)
	}

	return loadTest(testPath)
}

func loadPeripheralTest(pClass string, test string) (*template.Template, error) {
	testPath := path.Join(dirPeripherals, pClass, test)
	_, err := os.Stat(testPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("class '%s' is missing test '%s' test", pClass, test)
	} else if err != nil {
		return nil, fmt.Errorf("unable to open '%s' directory for class '%s': %w", test, pClass, err)
	}

	return loadTest(testPath)
}

func loadTest(testPath string) (*template.Template, error) {
	testBytes, err := ioutil.ReadFile(path.Join(testPath, "main.go"))
	if err != nil {
		return nil, fmt.Errorf("unable to read test '%s': %w", testPath, err)
	}

	tmpl, err := template.New("test").Parse(string(testBytes))
	if err != nil {
		return nil, fmt.Errorf("unable to parse 'main.go' as a Go template for test %s: %w", testPath, err)
	}

	return tmpl, nil
}

func runPeripheralTestForTarget(pClass string, name string, test *template.Template, target string, peripherals []string) ([]peripheralTestResult, error) {
	results := make([]peripheralTestResult, 0, len(peripherals))

	for _, peripheral := range peripherals {
		testResult, err := runTestForPeripheral(pClass, name, test, target, peripheral)
		if err != nil {
			return nil, fmt.Errorf("failed to run test for peripheral %s: %w", peripheral, err)
		}

		results = append(results, testResult)
	}

	return results, nil
}

func runTestForPeripheral(pClass string, name string, test *template.Template, target string, peripheral string) (peripheralTestResult, error) {
	output, err := buildFromTemplate(test, target, peripheral)
	return peripheralTestResult{
		Target:          target,
		PeripheralClass: pClass,
		Feature:         name,
		Peripheral:      peripheral,
		Passed:          err == nil,
		Output:          output,
	}, nil
}

func runTestForFeature(name string, test *template.Template, target string) (featureTestResult, error) {
	output, err := buildFromTemplate(test, target, "")
	return featureTestResult{
		Target:  target,
		Feature: name,
		Passed:  err == nil,
		Output:  output,
	}, nil
}
