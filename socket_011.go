package main

import (
	"fmt"
	"gopkg.in/qml.v0"
	"net"
	"os"
	"strings"
	"time"
)

type LogRecord struct {
	record string
	rType  string
	color  string
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

func (l *Logs) Color(index int) string {
	return l.list[index].color
}

type Bridge struct {
	started   bool
	ipName    string
	logs      *Logs
	listView  qml.Object
	inputText qml.Object
}

func (b *Bridge) SetupVars(root qml.Object) {
	b.listView = root.ObjectByName("logview")
	b.inputText = root.ObjectByName("inputtext1")
}

func (b *Bridge) ColorLog(log string, color string) {
	const layout = "Jan 2, 2006 at 3:04pm (MST)"
	strSplits := strings.Split(log, "\n")
	var firstTime bool = true
	for _, s := range strSplits {
		if len(s) > 0 {
			if firstTime {
				firstTime = false
				s = time.Now().Format(layout) + " : " + s
			}
			b.logs.Add(LogRecord{record: s, rType: "info", color: color})
		}
	}
}

func (b *Bridge) Log(log string) {
	b.ColorLog(log, "gold")
}

func (b *Bridge) ClientLog(log string) {
	b.ColorLog(log, "green")
}

func (b *Bridge) ErrorLog(log string) {
	b.logs.Add(LogRecord{record: log, rType: "error"})
}

func (b *Bridge) HandleClick() {
	b.ClientLog("Clicked")
	go b.UDPClient()
}

func checkSum(msg []byte) uint16 {
	sum := 0

	for n := 1; n < len(msg)-1; n += 2 {
		sum += int(msg[n])*256 + int(msg[n+1])
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var answer uint16 = uint16(^sum)
	return answer
}

func (b *Bridge) UDPClient() {
	host := "127.0.0.1"
	addr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		b.ClientLog(fmt.Sprintf("Resolution error: %s", err.Error()))
		return
	}

	conn, err := net.DialIP("ip:icmp", nil, addr)
	if err != nil {
		b.ClientLog(fmt.Sprintf("Error Dial: %s", err.Error()))
		return
	}

	var msg [512]byte
	msg[0] = 8  // echo
	msg[1] = 0  // code 0
	msg[2] = 0  // checksum fix later
	msg[3] = 0  // checksum fix later
	msg[4] = 0  // identifier[0]
	msg[5] = 13 // identifier[1]
	msg[6] = 0  // sequence[0]
	msg[7] = 37 // sequence[1]
	len := 8

	check := checkSum(msg[0:len])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 255)

	_, err = conn.Write(msg[0:len])
	if err != nil {
		b.ClientLog(fmt.Sprintf("Error send: %s", err.Error()))
		return
	}

	_, err = conn.Read(msg[0:])
	if err != nil {
		b.ClientLog(fmt.Sprintf("Error read: %s", err.Error()))
		return
	}

	b.ClientLog("Got Response")

	if msg[5] == 13 {
		b.ClientLog("identifier matches")
	}

	if msg[7] == 37 {
		b.ClientLog("Sequence matches")
	}
}

func main() {

	if err := run("socket_011.qml"); err != nil {
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

	bridge.SetupVars(win.Root())

	win.Show()
	win.Wait()
	return nil
}
