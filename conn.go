package main

import (
	"./access"
	"./problem"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"time"
	// "os/exec"
	// "math"
	"strings"
	"sync"
)

var connections = 0
var conn_lock = new(sync.Mutex)
var glob_lock = new(sync.Mutex)
var globalAccess *access.Access = access.NewAccess()
var globalSolving *access.Access = access.NewAccess()

var next_id = 1
var id_lock = new(sync.Mutex)

// var deadtime, _ = time.ParseDuration("2s")

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	difficulty      problem.Difficulty
	problems        []problem.Problem
	id              int
	access, solving *access.Access
}

type message struct {
	Opcode, SocketId    int
	Result, Query, Hash string
	Problems            []problem.Problem
	Difficulty          problem.Difficulty
}

func (m message) String() string {
	return fmt.Sprintf("Id: %d, Opcode: %d, Difficulty: %d, Query: %s", m.SocketId, m.Opcode, m.Difficulty, m.Query)
}

func (c *connection) reader() {
	conn_lock.Lock()
	connections++
	conn_lock.Unlock()
	// cmd := exec.Command("top", "-n", "0", "-stats", "cpu", "-l", "1")
	// // cmd.Start()
	// cmd.Wait()
	// o, err := cmd.Output()
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("Output: %s", o)
	conn_ok := true
	for conn_ok {
		var msg message
		err := websocket.JSON.Receive(c.ws, &msg)
		// fmt.Printf("Received: %v, ERR: %v\n", msg, err)
		if err != nil {
			break
		}
		log.Printf("RECEIVED: %v\n", msg)
		// c.ws.SetDeadline(time.Now().Add(deadtime))
		var response message
		if msg.Opcode == 0 {
			c.difficulty = problem.GetDifficulty(problem.Param{Local: *c.access, Global: *globalAccess, Cpu: CPU_STAT, LSolve: *c.solving, GSolve: *globalSolving}) // math.Max(CPU_LOAD, CPU_AVG)
			response.Difficulty = c.difficulty
			c.problems = problem.ConstructProblemSet(c.difficulty)
			response.Problems = c.problems
			c.access.Touch(1)
			c.solving.Caress()
			glob_lock.Lock()
			globalAccess.Touch(connections)
			glob_lock.Unlock()
			log.Printf("Local average: %v, Global Average: %v", c.access.ShortMean/1000000, globalAccess.ShortMean/1000000)

			response.Opcode = 1
			response.Query = msg.Query
			response.SocketId = c.id
		} else {
			if problem.Verify(c.problems, msg.Problems, c.difficulty) {
				diff := c.solving.Touch(1)
				glob_lock.Lock()
				globalSolving.AddTime(diff, connections)
				glob_lock.Unlock()

				response.Query = strings.Join([]string{"Your query, \"", msg.Query, "\" has been served since you solved the puzzle."}, "")
				before := time.Now().UTC().UnixNano()
				for i := int64(0); i < 250000000; i++ {
					//simulate some server load (~80 ms)
				}
				after := time.Now().UTC().UnixNano()
				log.Printf("Service took %.3f ms", float64(after-before)/float64(1000000))
			} else {
				response.Opcode = 255
				response.Result = "Incorrect hash!"
				response.Query = "Your query was ignored since you did not solve the puzzle."
				conn_ok = false
			}
		}
		log.Printf(" SENDING: %v\n", response)
		websocket.JSON.Send(c.ws, response)
	}
	conn_lock.Lock()
	connections--
	conn_lock.Unlock()
}

// func (c *connection) writer() {
// 	for message := range c.send {
// 		err := websocket.Message.Send(c.ws, message)
// 		if err != nil {
// 			break
// 		}
// 	}
// 	c.ws.Close()
// }

func wsHandler(ws *websocket.Conn) {
	id_lock.Lock()
	id := next_id
	next_id++
	id_lock.Unlock()
	log.Printf("Accepted connection from %s, assigning id %d\n", ws.Request().RemoteAddr, id)
	c := &connection{ws: ws, id: id, access: access.NewAccess(), solving: access.NewAccess()}

	// c.ws.SetDeadline(time.Now().Add(deadtime))
	//h.register <- c
	// defer func() { h.unregister <- c }()
	// go c.writer()
	c.reader()
	log.Printf("Client %d disconnected from server\n", id)

}
