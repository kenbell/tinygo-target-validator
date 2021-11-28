# TinyGo Target Validator

TinyGo Target Validator checks all known TinyGo target devices to see if their interface complies to the common target capabilities, and which optional devices they support.

## Usage
To build and use this tool:

```shell
go build .
./tinygo-target-validator test
```

To speed up testing for specific TinyGo targets, or specific peripheral classes use the `-target` and `-pclass` options, e.g.:
```shell
./tinygo-target-validator test -pclass spi -target pico
```

## Scope

This tool does not validate the implementation for a particular TinyGo target is correct - the tool is just intended to ensure that the expected interface is exposed.

## Output

The tool outputs a JSON document with a matrix of all targets and their peripherals.

## Implementation

The tool works by trying to compile a suite of test code against each target.  If the code compiles, the target is deemed 'compatible'.  If the compilation fails, the target is deemed 'incompatible'.

## Structure

The tests are structured under the `tests` folder.  There are two classes of test:

* **Feature tests** which test the presence of expected functions / types / declarations
* **Peripheral tests** which test MCU peripherals meet expected interfaces



### Feature Tests

These tests are present under:
```
tests/features/<feature>
```

The feature tests each consist of a `main.go`.  The code itself does not need to do anything other than ensure that it compiles without error for a compliant device and fails to compile for non-compliant devices.

An example test is `pin-interrupt`:
```go
package main

import "machine"

func main() {
	machine.Pin(0).SetInterrupt(
		machine.PinRising,
		func(machine.Pin) {})
}
```

This test verifies if the Pin type supports the `SetInterrupt` function with the correct signature.

### Peripheral Tests

These tests have this folder structure:
```
tests/peripherals/<peripheral class>/<capability>
```

The peripheral class is the type of peripheral, eg `uart` or `spi`.  The capability is either `core` indicating the bare minimum required functionality or can be any meaningful name, indicating an optional capability.  An example is:
```
tests/peripherals/i2c/bus-speed-control
```

This would be tests to validate whether an i2c peripheral supports the interface for controlling the I2C bus speed.

The individual tests are basic Go apps that are run through the Go template language and compiled.  A prototypical example for peripherals is to ensure the peripheral meets an interface, such as this example for I2C:
```go
package main

import (
	"machine"
)

type i2c interface {
	Configure(machine.I2CConfig) error
	Tx(x uint16, w []byte, r []byte) error
}

func main() {
	_ = i2c(machine.{{.Peripheral}})
}
```

The code itself does not need to do anything other than ensure that it compiles without error when a peripheral meets the required interface and fails to compile if it doesn't meet the required interface.