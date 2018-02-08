package main

import (
	"github.com/joyent/triton-go/cmd/internal/console_writer"
	"github.com/joyent/triton-go/cmd/triton/cmd"
)

func main() {
	defer func() {
		p := console_writer.GetTerminal()
		p.Wait()
	}()

	cmd.Execute()
}
