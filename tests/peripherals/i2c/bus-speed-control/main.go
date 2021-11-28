package main

import (
	"machine"
)

type i2c interface {
	Configure(machine.I2CConfig) error
	ReadRegister(addr uint8, r uint8, buf []byte) error
	WriteRegister(addr uint8, r uint8, buf []byte) error
	Tx(x uint16, w []byte, r []byte) error
}

func main() {
	p := i2c(machine.{{.Peripheral}})
	p.Configure(machine.I2CConfig{Frequency:400000})
}
