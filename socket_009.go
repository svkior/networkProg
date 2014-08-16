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
	started  bool
	ipName   string
	logs     *Logs
	listView qml.Object
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
	b.logs.Add(LogRecord{record: time.Now().String(), rType: "info"})
	for _, s := range strSplits {
		if len(s) > 0 {
			b.logs.Add(LogRecord{record: s, rType: "info"})
		}
	}
	//b.listView.Set("currentIndex", b.logs.Len-1)
}

func (b *Bridge) ErrorLog(log string) {
	b.logs.Add(LogRecord{record: log, rType: "error"})
}

func (b *Bridge) EchoServer() {
	service := ":1201"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		b.Log(fmt.Sprintf("Error get tcpAddr: %s", err.Error()))
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		b.Log(fmt.Sprintf("Error listen tcpAddr: %s", err.Error()))
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go b.handleClient(conn)
	}
}

func (b *Bridge) handleClient(conn net.Conn) {

	defer conn.Close()

	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		b.Log(string(buf[0:n]))
		_, err = conn.Write(buf[0:n])
		if err != nil {
			return
		}
	}
}

func main() {

	if err := run("socket_009.qml"); err != nil {
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

	bridge.listView = win.Root().ObjectByName("logview")

	win.Show()
	go bridge.EchoServer()
	win.Wait()
	return nil
}
