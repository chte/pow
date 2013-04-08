package dataminer

import (
	"bufio"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

var controlTempl = template.Must(template.ParseFiles("./dataminer/control.html"))
var logging bool = false

func controlHTMLHandler(c http.ResponseWriter, req *http.Request) {
	controlTempl.Execute(c, req.Host)
}

// var datacollection = {    "Behaviour": behaviour_type, 
//         "RequestingTime": requestingTime,
//         "SolvingTime": solvingTime,
//         "GrantingTime": grantingTime};

type sample struct {
	Behaviour                                 string
	RequestingTime, SolvingTime, GrantingTime int
}

func (s *sample) String() string {
	// return fmt.Sprintf("%d & %d & %d & %d", s.RequestingTime, s.SolvingTime, s.GrantingTime, s.RequestingTime+s.SolvingTime+s.GrantingTime)
	return fmt.Sprintf("%d	%d	%d	%d", s.RequestingTime, s.SolvingTime, s.GrantingTime, s.RequestingTime+s.SolvingTime+s.GrantingTime)

}

type ctrlMsg struct {
	Command int
	Name    string
}

type hub struct {
	control chan ctrlMsg
	log     chan sample
}

var h = hub{
	control: make(chan ctrlMsg),
	log:     make(chan sample, 100),
}

type Data struct {
	series map[string]*dataseries
}

func (d *Data) getSeries(s string) *dataseries {
	if d.series[s] == nil {
		d.series[s] = NewDataSeries()
	}
	return d.series[s]
}

type dataseries struct {
	data []sample
}

func NewDataSeries() *dataseries {
	d := new(dataseries)
	d.data = make([]sample, 0, 100)
	return d
}
func (d *dataseries) add(s sample) {
	d.data = append(d.data, s)
}
func NewData() *Data {
	var d Data
	d.series = make(map[string]*dataseries)
	return &d
}

var data *Data

func getFormula(formula, who string, n int) string {
	return "=" + formula + "(" + who + "2:" + who + strconv.Itoa(n+1) + ")"
}

func getConfidence(alpha string, who string, index_of_avg int, n int) string {
	return "=CONFIDENCE.T(" + alpha + "," + who + strconv.Itoa(index_of_avg) + "," + strconv.Itoa(n) + ")"
}

func (ds *dataseries) WriteTo(filename string) {
	// f, err := os.Create(filename)
	// wd, _ := os.Getwd()
	fil, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	// close fil on exit and check for its returned error
	defer func() {
		if fil.Close() != nil {
			panic(err)
		}
	}()
	w := bufio.NewWriter(fil)
	fmt.Fprintf(w, "%s	%s	%s	%s\n", "request", "solve", "grant", "service")
	for _, s := range ds.data {
		fmt.Fprintf(w, "%v\n", &s)
	}
	n := len(ds.data)
	formula := "STDEV.P"
	fmt.Fprintf(w, "%s	%s	%s	%s	STDDEV\n", getFormula(formula, "A", n), getFormula(formula, "B", n), getFormula(formula, "C", n), getFormula(formula, "D", n))
	formula = "AVERAGE"
	fmt.Fprintf(w, "%s	%s	%s	%s	AVERAGE\n", getFormula(formula, "A", n), getFormula(formula, "B", n), getFormula(formula, "C", n), getFormula(formula, "D", n))
	fmt.Fprintf(w, "%s	%s	%s	%s	CONFIDENCE 0.05\n", getConfidence("0.05", "A", n+2, n), getConfidence("0.05", "B", n+2, n), getConfidence("0.05", "C", n+2, n), getConfidence("0.05", "D", n+2, n))
	fmt.Fprintf(w, "%s	%s	%s	%s	CONFIDENCE 0.01\n", getConfidence("0.01", "A", n+2, n), getConfidence("0.01", "B", n+2, n), getConfidence("0.01", "C", n+2, n), getConfidence("0.01", "D", n+2, n))

	w.Flush()
}

func (h *hub) run() {
	for {
		select {
		case cmd := <-h.control:
			switch cmd.Command {
			case 0:
				log.Printf("Start recording\n")
				data = NewData()
				logging = true
				break
			case 1:
				log.Printf("Stop Recording series %d, saving to file\n", cmd.Name)
				t := time.Now()
				for k, v := range data.series {
					v.WriteTo("./data/" + t.Format("15_04_05") + "-" + cmd.Name + "-" + k + ".pow_data")
				}
				logging = false
				break
			}
			break
		case smp := <-h.log:
			// log.Printf("Received sample %v\n", smp)
			if logging {
				data.getSeries(smp.Behaviour).add(smp)

			}
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
