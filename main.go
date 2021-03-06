package main
  
import (
	"./dataminer"
	"./problem"
	"code.google.com/p/go.net/websocket"
	"flag"
	//"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	//"runtime"
	"text/template"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")
var homeTempl = template.Must(template.ParseFiles("./home.html"))
var attackTempl = template.Must(template.ParseFiles("./attack.html"))
var jsTempl = template.Must(template.ParseFiles("./pow.js"))
var jsAtkTempl = template.Must(template.ParseFiles("./pow_atk.js"))
var jsAtkTaskTempl = template.Must(template.ParseFiles("./attacktask.js"))
var monitorTempl = template.Must(template.ParseFiles("./monitor.html"))

var serveStaticMap = make(map[string]*template.Template)

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, problem.BaseDifficulty)
}
func attackHandler(c http.ResponseWriter, req *http.Request) {
	attackTempl.Execute(c, problem.BaseDifficulty)
}
func ipHandler(c http.ResponseWriter, req *http.Request) {
	// c.WriteHeader()
	c.Write([]byte(req.RemoteAddr))
}

// func jsAtkHandler(c http.ResponseWriter, req *http.Request) {
// 	jsAtkTempl.Execute(c, req.Host)
// }

// func jsAtkTaskHandler(c http.ResponseWriter, req *http.Request) {
// 	jsAtkTaskTempl.Execute(c, req.Host)
// }

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
	//cpu := runtime.GOMAXPROCS(runtime.NumCPU())
	//log.Printf("PoW server started on %s server at %s, threads: %d -> %d\n", runtime.GOOS, *addr, cpu, runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())
	go h.run()
	dataminer.Run()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		defer func() {
			//fmt.Println("hallo")
		}()
		/*for sig := range c {
			//fmt.Printf("sig: %v", sig)
			os.Exit(0)
			// return
			// sig is a ^C, handle it
		}*/
	}()
	http.HandleFunc("/", homeHandler)
	//http.HandleFunc("/pow.js", jsHandl	er)
	http.HandleFunc("/attack", attackHandler)
	// http.HandleFunc("/pow_atk.js", jsAtkHandler)
	// http.HandleFunc("/attacktask.js", jsAtkTaskHandler)

	// http.HandleFunc("/pow.js", jsHandler)
	serveStatic("sha256.js")
	serveStatic("jquery.js")
	serveStatic("pow.js")
	serveStatic("highcharts.js")
	serveStatic("highcharts_theme.js")
	serveStatic("exporting.js")
	serveStatic("pow_atk.js")
	serveStatic("attacktask.js")
	http.HandleFunc("/monitor", monitorHTMLHandler)
	http.HandleFunc("/ip", ipHandler)
	http.Handle("/ws", websocket.Handler(wsHandler))
	http.Handle("/monitor_ws", websocket.Handler(monitorHandler))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
