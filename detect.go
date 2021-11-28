package main

// This file contains the code to detect TinyGo peripherals.
//
// Currently this is very slow, with a test compilation performed
// to detect each peripheral.  Peripherals are assumed to be named
// like this: SPIn
//
// That-is, the peripheral class (SPI in this case) followed by an
// integer starting at 0.
//
// Currently the code tries all peripherals between the values of
// 0 and 7.

import (
	"fmt"
	"html/template"
	"strings"
)

const peripheralDetect = `package main

import "machine"

func main() {
	_ = machine.{{.Peripheral}}
}`

func detectPeripherals(target string, pClass string) ([]string, error) {
	peripherals := make([]string, 0, 8)

	tmpl, err := template.New("test").Parse(peripheralDetect)
	if err != nil {
		return nil, fmt.Errorf("unable to parse peripheralDetect template, %w", err)
	}

	for i := 0; i < 8; i++ {
		pName := fmt.Sprintf("%s%d", strings.ToUpper(pClass), i)

		// If the build fails = no peripheral with this name
		_, err = buildFromTemplate(tmpl, target, pName)
		if err == nil {
			peripherals = append(peripherals, pName)
		}
	}

	return peripherals, nil
}
