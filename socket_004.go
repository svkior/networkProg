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
	addrs, err := net.LookupHost(name)

	if err != nil {
		resultIP.Set("text", fmt.Sprint("Resolution error:", err.Error()))
	} else {
		msg := ""
		for i, s := range addrs {
			msg += fmt.Sprintf("%s (%d) : %s\n", name, i, s)
		}
		resultIP.Set("text", msg)
		edtIP.Call("selectAll")
	}
}

func main() {

	if err := run("socket_004.qml"); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run(qmlName string) error {
	qml.Init(nil)
	engine := qml.NewEngine()

	component, err := engine.LoadFile(qmlName)
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
