package dataminer

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
	"text/template"
	"time"
)

var controlTempl = template.Must(template.ParseFiles("./dataminer/control.html"))

func controlHTMLHandler(c http.ResponseWriter, req *http.Request) {
	controlTempl.Execute(c, req.Host)
}

type sample struct {
	Tid time.Time
}

type ctrlMsg struct {
	Command int
}

type hub struct {
	control chan ctrlMsg
	log     chan sample
}

var h = hub{
	control: make(chan ctrlMsg),
	log:     make(chan sample, 100),
}

func (h *hub) run() {
	for {
		select {
		case cmd := <-h.control:
			switch cmd.Command {
			case 0:
				log.Printf("Start recording\n")
				break
			case 1:
				log.Printf("Stop Recording save to file\n")
				break
			}
			break
		case smp := <-h.log:
			log.Printf("Received sample %v\n", smp)
		}
	}
}

func controlHandler(ws *websocket.Conn) {
	var msg ctrlMsg
	for {
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			break
		}
		h.control <- msg
	}
	ws.Close()
}
func dataHandler(ws *websocket.Conn) {
	var smp sample
	for {
		err := websocket.JSON.Receive(ws, &smp)
		if err != nil {
			break
		}
		h.log <- smp
	}
	ws.Close()
}

func Run() {
	http.HandleFunc("/data_ctrl", controlHTMLHandler)
	http.Handle("/dataminer_ctrl", websocket.Handler(controlHandler))
	http.Handle("/dataminer", websocket.Handler(dataHandler))

	go h.run()
}
