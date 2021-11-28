package main

import "machine"

func main() {
	machine.Pin(0).SetInterrupt(
		machine.PinRising,
		func(machine.Pin) {})
}
