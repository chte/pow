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
var monitorTempl = template.Must(template.ParseFiles("./monitor.html"))

var serveStaticMap = make(map[string]*template.Template)

type homeParam struct {
	Difficulty int
	Problems   int
}

var defaultParam = homeParam{3, 16}

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, defaultParam)
}

// func jsHandler(c http.ResponseWriter, req *http.Request) {
// 	jsTempl.Execute(c, req.Host)
// }
func monitorHTMLHandler(c http.ResponseWriter, req *http.Request) {
	monitorTempl.Execute(c, req.Host)
}
func serveStatic(file string) {
	serveStaticMap[file] = template.Must(template.ParseFiles("./" + file))
	http.HandleFunc("/"+file, func(c http.ResponseWriter, req *http.Request) {
		serveStaticMap[file].Execute(c, req.Host)
	})
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())
	go h.run()
	http.HandleFunc("/", homeHandler)
	// http.HandleFunc("/pow.js", jsHandler)
	serveStatic("pow.js")
	serveStatic("highcharts.js")
	serveStatic("highcharts_theme.js")
	serveStatic("exporting.js")
	http.HandleFunc("/monitor", monitorHTMLHandler)
	http.Handle("/ws", websocket.Handler(wsHandler))
	http.Handle("/monitor_ws", websocket.Handler(monitorHandler))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
