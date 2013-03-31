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
var attackTempl = template.Must(template.ParseFiles("./attack.html"))
var jsTempl = template.Must(template.ParseFiles("./pow.js"))
var jsAtkTempl = template.Must(template.ParseFiles("./pow_atk.js"))
var jsAtkTaskTempl = template.Must(template.ParseFiles("./attacktask.js"))

type homeParam struct {
	Difficulty int
	Problems   int
}

var defaultParam = &homeParam{3, 16}

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, defaultParam)
}
func attackHandler(c http.ResponseWriter, req *http.Request) {
	attackTempl.Execute(c, req.Host)
}

func jsAtkHandler(c http.ResponseWriter, req *http.Request) {
	jsAtkTempl.Execute(c, req.Host)
}

func jsAtkTaskHandler(c http.ResponseWriter, req *http.Request) {
	jsAtkTaskTempl.Execute(c, req.Host)
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
	http.HandleFunc("/attack", attackHandler)
	http.HandleFunc("/pow_atk.js", jsAtkHandler)
	http.HandleFunc("/attacktask.js", jsAtkTaskHandler)

	http.Handle("/ws", websocket.Handler(wsHandler))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
