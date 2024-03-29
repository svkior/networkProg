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
	logs    *Logs
}

type LogRecord struct {
	record string
	rType  string
}

type Logs struct {
	list []LogRecord
	Len  int
}

func (l *Logs) Add(r LogRecord) {
	l.list = append(l.list, r)
	l.Len = len(l.list)
	qml.Changed(l, &l.Len)
}

func (l *Logs) Type(index int) string {
	return l.list[index].rType
}

func (l *Logs) Record(index int) string {
	return l.list[index].record
}

func (b *Bridge) Log(log string) {
	b.logs.Add(LogRecord{record: log, rType: "info"})
}

func (b *Bridge) ErrorLog(log string) {
	b.logs.Add(LogRecord{record: log, rType: "error"})
}

func (b *Bridge) HandleClick(edtIP qml.Object) {
	name := edtIP.String("text")
	addrs, err := net.LookupHost(name)

	if err != nil {
		b.Log(fmt.Sprintf("Resolution error: (%s)", err.Error()))
	} else {
		for i, s := range addrs {
			b.Log(fmt.Sprintf("%s (%d) : %s", name, i, s))
		}
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

	logs := Logs{}

	bridge := Bridge{
		started: false,
		ipName:  "127.0.0.1",
		logs:    &logs,
	}

	context := engine.Context()
	context.SetVar("bridge", &bridge)
	context.SetVar("logs", &logs)

	win := component.CreateWindow(nil)

	win.Show()
	win.Wait()
	return nil
}
