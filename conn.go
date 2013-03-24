package main

import (
	"code.google.com/p/go.net/websocket"
	"crypto/sha256"
	// "encoding/json"
	"encoding/hex"
	"fmt"
	"hash"
	"math/rand"
	"strconv"
	"strings"
)

const NUMBER_OF_PROBLEMS int = 16

type connection struct {
	// The websocket connection.
	ws *websocket.Conn
	// json_out         *json.Encoder
	// json_in          *json.Decoder
	ha         hash.Hash
	difficulty int
	problems   []problem

	// Buffered channel of outbound messages.
	// send chan string
}

type problem struct {
	Seed, Solution int
}

func Newproblem() problem {
	return problem{Seed: int(rand.Uint32())}
}

type message struct {
	Opcode, Difficulty  int
	Result, Query, Hash string
	Problems            []problem
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
	fmt.Printf("Client connected\n")
	for {
		var msg message
		err := websocket.JSON.Receive(c.ws, &msg)
		// fmt.Printf("test: %v, ERR: %v\n", msg, err)
		if err != nil {
			break
		}
		fmt.Printf("Got query: %v\n", msg)
		var response message
		if msg.Opcode == 0 {
			c.difficulty, response.Difficulty = defaultHomeParam.Difficulty, defaultHomeParam.Difficulty
			response.Problems = make([]problem, NUMBER_OF_PROBLEMS)
			for i, _ := range response.Problems {
				response.Problems[i] = Newproblem()
				c.problems[i] = response.Problems[i]
			}
			response.Opcode = 1
			response.Query = msg.Query
			// c.json_out.Encode(response)
		} else {
			ok := true
			var sha string
			for i := 0; i < NUMBER_OF_PROBLEMS; i++ {
				c.ha.Reset()
				c.ha.Write([]byte(strconv.Itoa(msg.Problems[i].Solution)))
				c.ha.Write([]byte(strconv.Itoa(c.problems[i].Seed)))
				// has := make([]byte, c.ha.Size())
				// has = c.ha.Sum(has)
				sha = hex.EncodeToString(c.ha.Sum(nil))
				fmt.Printf("Response solution: %v\n Calc Solution: %v\n", msg.Problems[i].Solution, sha)
				if init_zeroes(sha) < c.difficulty {
					ok = false
					break
				}
			}
			if ok {
				// response.Result = strconv.Itoa(msg.Solution)
				response.Query = strings.Join([]string{"Your query, \"", msg.Query, "\" has been served since you solved the puzzle."}, "")
				response.Hash = sha //string(has)
			} else {
				// response.Result = strings.Join([]string{"Incorrect hash!: ", strconv.Itoa(msg.Problems.Solution)}, "")
				response.Query = "Your query was ignored since you did not solve the puzzle."
			}
		}
		fmt.Printf("Sending response %v\n", response)
		websocket.JSON.Send(c.ws, response)
	}
	fmt.Printf("Client Disconnected\n")
	c.ws.Close()
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
	c := &connection{ha: sha256.New(), ws: ws, problems: make([]problem, NUMBER_OF_PROBLEMS)}
	//h.register <- c
	// defer func() { h.unregister <- c }()
	// go c.writer()
	c.reader()
}
