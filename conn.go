package main

import (
	"code.google.com/p/go.net/websocket"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"log"
	"math/rand"
	// "os/exec"
	"strconv"
	"strings"
	"sync"
)

var NUMBER_OF_PROBLEMS = defaultParam.Problems
var connections = 0
var conn_lock = new(sync.Mutex)

var next_id = 0
var id_lock = new(sync.Mutex)

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	ha         hash.Hash
	difficulty int
	problems   []problem
	id         int
}

type problem struct {
	Seed, Solution int
}

func Newproblem() problem {
	return problem{Seed: int(rand.Uint32())}
}

type message struct {
	Opcode, Difficulty, ConnId int
	Result, Query, Hash        string
	Problems                   []problem
}

func init_zeroes(s string) (num int) {
	for _, c := range s {
		if c != '0' {
			return
		}
		num++
	}
	return
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
	for {
		var msg message
		err := websocket.JSON.Receive(c.ws, &msg)
		// fmt.Printf("Received: %v, ERR: %v\n", msg, err)
		if err != nil {
			break
		}
		log.Printf("Got query: %v\n", msg)
		var response message
		if msg.Opcode == 0 {
			c.difficulty, response.Difficulty = defaultParam.Difficulty, defaultParam.Difficulty
			response.Problems = make([]problem, NUMBER_OF_PROBLEMS)
			for i, _ := range response.Problems {
				response.Problems[i] = Newproblem()
				c.problems[i] = response.Problems[i]
			}
			response.Opcode = 1
			response.Query = msg.Query
			response.ConnId = c.id
		} else {
			ok := true
			var sha string
			for i := 0; i < NUMBER_OF_PROBLEMS; i++ {
				c.ha.Reset()
				c.ha.Write([]byte(strconv.Itoa(msg.Problems[i].Solution)))
				c.ha.Write([]byte(strconv.Itoa(c.problems[i].Seed)))
				sha = hex.EncodeToString(c.ha.Sum(nil))
				// fmt.Printf("Response solution: %v\n Calc Solution: %v\n", msg.Problems[i].Solution, sha)
				if init_zeroes(sha) < c.difficulty {
					ok = false
					break
				}
			}
			if ok {
				response.Query = strings.Join([]string{"Your query, \"", msg.Query, "\" has been served since you solved the puzzle."}, "")
				response.Hash = sha
				for i := 0; i < 10000000; i++ {
					//simulate some server load
				}
			} else {
				response.Result = "Incorrect hash!"
				response.Query = "Your query was ignored since you did not solve the puzzle."
			}
		}
		log.Printf("Sending response to %d: %v\n", c.id, response)
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
	log.Printf("Accepted connection from %s, assigning id %d\n", ws.LocalAddr(), id)
	c := &connection{ha: sha256.New(), ws: ws, problems: make([]problem, NUMBER_OF_PROBLEMS), id: id}
	//h.register <- c
	// defer func() { h.unregister <- c }()
	// go c.writer()
	c.reader()
	log.Printf("Client %d disconnected from server\n", id)

}
