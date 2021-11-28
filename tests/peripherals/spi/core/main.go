package main

import (
	"machine"
)

type spi interface {
	Configure(config machine.SPIConfig) error
	Tx(w, r []byte) (err error)
	Transfer(w byte) (byte, error)
}

func main() {
	p := spi(machine.{{.Peripheral}})
	p.Configure(machine.SPIConfig{Frequency: uint32(0), LSBFirst: true, Mode: uint8(0)})
}
