package main

import (
	"fmt"
	"gopkg.in/qml.v0"
	"net"
	"os"
	"strings"
	"time"
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

func main() {

	if err := run("socket_007.qml"); err != nil {
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
	go func() {
		service := ":1200"
		tcpAddr, err := net.ResolveTCPAddr("tcp", service)
		if err != nil {
			bridge.Log(fmt.Sprintf("Error get tcpAddr: %s", err.Error()))
			return
		}
		listener, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			bridge.Log(fmt.Sprintf("Error listen tcpAddr: %s", err.Error()))
			return
		}

		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}

			daytime := time.Now().String()
			remoteHost := conn.RemoteAddr().String()
			bridge.Log(fmt.Sprintf("Send %s to %s", daytime, remoteHost))
			conn.Write([]byte(daytime))
			conn.Close()
		}
	}()
	win.Wait()
	return nil
}
