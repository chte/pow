package main

import (
	"bufio"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
)

var CPU_LOAD float64
var CPU_AVG float64 = 40

const CPU_ALPHA = 0.01

type hub struct {
	// Registered connections.
	connections map[*subscriber]bool

	// Inbound messages from the connections.
	broadcast chan information

	// Register requests from the connections.
	register chan *subscriber

	// Unregister requests from connections.
	unregister chan *subscriber
}

var h = hub{
	broadcast:   make(chan information),
	register:    make(chan *subscriber),
	unregister:  make(chan *subscriber),
	connections: make(map[*subscriber]bool),
}

type information struct {
	Cpu_user, Cpu_system              float64
	Monitoring, Users                 int
	ShortAverageTime, LongAverageTime int64
}
type subscriber struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan information
}

//Collect cpu info with: top -n 0 -stats cpu -l 0
//cmd := exec.Command("top", "-n", "0", "-stats", "cpu", "-l", "0")
func collect(ch chan information) {
	//Collect cpu info with: top -n 0 -stats cpu -l 0
	//cmd := exec.Command("top", "-n", "0", "-stats", "cpu", "-l", "0") //Mac
	//Ubuntu

	mac := runtime.GOOS == "darwin"
	var cmd *exec.Cmd
	if mac {
		cmd = exec.Command("top", "-n", "0", "-stats", "cpu", "-l", "0") //Mac
	} else {
		//cmd = exec.Command("sar", "-u", "1")
		cmd = exec.Command("top", "-b", "-d", "1")
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	// stderr, err := cmd.StderrPipe()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	rd := bufio.NewReader(stdout)
	var exp *regexp.Regexp
	if mac {
		exp, _ = regexp.Compile("CPU usage: (.*?)% user, (.*?)% sys, (.*?)% idle") //Mac
	} else {

		//exp, _ = regexp.Compile("all\\s+(\\S+)\\s+\\S+\\s+(\\S+)")
		exp, _ = regexp.Compile("\\s(\\S*?)[ %]us,  (.*?)[ %]sy")
	}
	for {
		line, _, _ := rd.ReadLine()
		//fmt.Printf("%s\n", line)

		catch := exp.FindStringSubmatch(string(line))
		if catch != nil {
			//log.Printf("%s\n", catch[0])

			user, _ := strconv.ParseFloat(catch[1], 64)
			system, _ := strconv.ParseFloat(catch[2], 64)
			//log.Printf("%f(%s), %f(%s)\n", user, catch[1], system, catch[2])
			//log.Printf("%d", globalAccess.ShortMean)
			ch <- information{Cpu_user: user, Cpu_system: system, Monitoring: len(h.connections), Users: connections, ShortAverageTime: globalAccess.ShortMean, LongAverageTime: globalAccess.LongMean}
			CPU_LOAD = user
			CPU_AVG = CPU_ALPHA*CPU_LOAD + (1-CPU_ALPHA)*CPU_AVG
		}

	}

}
func (h *hub) run() {
	go collect(h.broadcast)
	for {
		select {
		case ws := <-h.register:
			h.connections[ws] = true
		case ws := <-h.unregister:
			delete(h.connections, ws)
			close(ws.send)
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					delete(h.connections, c)
					close(c.send)
					go c.ws.Close()
				}
			}
		}
	}
}

func (c *subscriber) writer() {
	for message := range c.send {
		err := websocket.JSON.Send(c.ws, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

func monitorHandler(ws *websocket.Conn) {
	log.Printf("Accepted monitor connection from %s\n", ws.RemoteAddr())
	c := &subscriber{send: make(chan information, 256), ws: ws}
	h.register <- c
	defer func() { h.unregister <- c }()
	c.writer()
}
