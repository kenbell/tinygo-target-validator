package main

import (
	"machine"
)

type pwm interface {
	Configure(config machine.PWMConfig) error
	SetPeriod(period uint64) error
	Top() uint32
	Channel(pin machine.Pin) (uint8, error)
	SetInverting(channel uint8, inverting bool)
	Set(channel uint8, value uint32)
}

func main() {
	p := pwm(machine.{{.Peripheral}})
	p.Configure(machine.PWMConfig{Period: uint64(0)})
}
