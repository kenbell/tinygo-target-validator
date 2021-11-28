package main

import (
	"io"
	"machine"
)

type uart interface {
	io.Reader
	io.Writer

	Buffered() int
}

func main() {
	_ = uart(machine.{{.Peripheral}})
}
