package main

import "machine"

func main() {
	var freq uint32
	freq = machine.CPUFrequency()
	_ = freq
}
