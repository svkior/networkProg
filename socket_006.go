package main

import (
	"fmt"
	"gopkg.in/qml.v0"
	"io/ioutil"
	"net"
	"os"
	"strings"
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
	strSplits := strings.Split(log, "\n")
	for _, s := range strSplits {
		if len(s) > 0 {
			b.logs.Add(LogRecord{record: s, rType: "info"})
		}
	}
}

func (b *Bridge) ErrorLog(log string) {
	b.logs.Add(LogRecord{record: log, rType: "error"})
}

func (b *Bridge) HandleClick(inpService qml.Object) {
	serivce := inpService.String("text")

	tcpAddr, err := net.ResolveTCPAddr("tcp4", serivce)
	if err != nil {
		b.Log(fmt.Sprintf("Error: (%s)", err.Error()))
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		b.Log(fmt.Sprintf("Error: (%s)", err.Error()))
		return
	}

	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	if err != nil {
		b.Log(fmt.Sprintf("Error: (%s)", err.Error()))
		conn.Close()
		return
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		b.Log(fmt.Sprintf("Error: (%s)", err.Error()))
		conn.Close()
		return
	}

	b.Log(string(result))
	inpService.Call("selectAll")
}

func main() {

	if err := run("socket_006.qml"); err != nil {
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
