package main

import (
	"fmt"
	"gopkg.in/qml.v0"
	"net"
	"os"
)

type Bridge struct {
	started bool
	ipName  string
}

func (b *Bridge) HandleClick(edtIP qml.Object, resultIP qml.Object) {
	name := edtIP.String("text")

	addr, err := net.ResolveIPAddr("ip", name)
	if err != nil {
		resultIP.Set("text", fmt.Sprint("Resolution error:", err.Error()))
	} else {
		resultIP.Set("text", addr.String())
		edtIP.Call("selectAll")
	}
}

func main() {

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	qml.Init(nil)
	engine := qml.NewEngine()

	component, err := engine.LoadFile("socket_003.qml")
	if err != nil {
		return err
	}

	bridge := Bridge{
		started: false,
		ipName:  "127.0.0.1",
	}
	context := engine.Context()
	context.SetVar("bridge", &bridge)

	win := component.CreateWindow(nil)

	win.Show()
	win.Wait()
	return nil
}
