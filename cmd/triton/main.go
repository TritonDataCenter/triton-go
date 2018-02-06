package main

import (
	"github.com/joyent/triton-go/cmd/internal/console_writer"
	"github.com/joyent/triton-go/cmd/triton/cmd"
)

func main() {
	p := console_writer.GetTerminal()
	defer p.Close()

	cmd.Execute()
}
