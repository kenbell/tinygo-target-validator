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
