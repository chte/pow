package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")
var homeTempl = template.Must(template.ParseFiles("./home.html"))
var jsTempl = template.Must(template.ParseFiles("./pow.js"))

type homeParam struct {
	Difficulty int
	Problems   int
}

var defaultParam = &homeParam{2, 256}

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, defaultParam)
}
func jsHandler(c http.ResponseWriter, req *http.Request) {
	jsTempl.Execute(c, req.Host)
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())
	// go h.run()
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/pow.js", jsHandler)
	http.Handle("/ws", websocket.Handler(wsHandler))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
