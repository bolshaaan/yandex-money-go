package main

import (
	"net/http"
	"regexp"
	"html/template"
	"io/ioutil"
	"time"
	"./ya-money-go"
	"fmt"
	"encoding/json"
	"log"
	"strconv"
	"github.com/GeertJohan/go.rice"
)

type Page struct {
	Title string
	Sum float64
	Operations string
	Body []byte
}

var file_name = "/home/alexander/token_file"
var validPath = regexp.MustCompile("^/(operations)/$")
var validDataPath = regexp.MustCompile("^/(data)/(in|out)$")

const template_path = "templates"

func operationsHandler(w http.ResponseWriter, r *http.Request) {

	p := &Page{}
	renderTemplate(w, "operations", p)
}

func makeHandler (fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)

		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func renderTemplate( w http.ResponseWriter, tmpl string, p *Page ) {
	t, _ := template.ParseFiles( template_path +  "/"  + tmpl + ".html")
	t.Execute(w, p)
}

func data_handler (w http.ResponseWriter, r *http.Request) {


	token, err := ioutil.ReadFile(file_name)

	if err != nil {
		panic(err)
	}

	ya, _ := yamoney.NewYaMoney(string(token))

	// Start period
	from, _ := time.Parse( "2006-Jan-02 00:00:00", "2016-Aug-01 00:00:00" )
	till := from.AddDate(0, 1, 0)

	operations := ya.OperationHistory( "from=" + from.Format(time.RFC3339) + "&till=" + till.Format(time.RFC3339) )

	m := validDataPath.FindStringSubmatch(r.URL.Path)

	if m == nil {
		http.NotFound(w,r)
		return
	}

	var data struct {
		In, Out []interface{}
		Labels map[string]string
		Sum    map[string]float64
		AggregatedOut []interface{}
		Total float64
	}
	var sum float64

	data.Labels = make(map[string]string)
	data.Sum = make(map[string]float64)

	dur, err := time.ParseDuration("24h")
	check(err)

	//run around and make sum

	for i := len(operations)-1; i>=0; i-- {

		op := operations[i]

		if op.Direction != "out" {
			continue
		}

		fmt.Println( op.DateTime, " : ", op.Title, " : ",  op.Amount)
		sum += op.Amount

		curt, _ := time.Parse(time.RFC3339, op.DateTime)

		var rr *[]interface{}
		if op.Direction == "in" {
			rr = &data.In
		} else {
			rr = &data.Out
		}

		kk := curt.Round(dur).Unix() * 1000
		k := strconv.FormatInt(curt.Round(dur).Unix() * 1000, 10)

		data.Labels[k] = fmt.Sprintf(
			//`%s <tr><td><span style="color:green">%s</span></td><td>%f</td></tr>`, data.Labels[k],
			`%s  ★ %06.2f ₽ <span style="color: #78a2b7">%s</span><br/>`, data.Labels[k],
			op.Amount, op.Title )

		*rr = append( *rr, []interface{}{ kk, op.Amount })

		data.Sum[ k ] += op.Amount
		data.Total += op.Amount
	}

	//data.AggregatedOut = make([]interface{}, len(data.Sum))
	for i := len(operations)-1; i>=0; i-- {
		op := operations[i]

		if op.Direction != "out" {
			continue
		}

		curt, _ := time.Parse(time.RFC3339, op.DateTime)
		k := strconv.FormatInt(curt.Round(dur).Unix() * 1000, 10)


		data.AggregatedOut = append(data.AggregatedOut, []interface{}{ curt.Round(dur).Unix() * 1000, data.Sum[k] })
	}

	for k := range data.Labels {

		data.Labels[k] = fmt.Sprintf(
			`<br/><span style="background-color:yellow;color:red;font-size:+10">%f</span> <br/>%s`,
			data.Sum[k], data.Labels[k] )
	}



	log.Println( "Labels: ", data.Labels )

	jj, err := json.Marshal(data)
	check(err)

	log.Println( string(jj) )

	fmt.Fprint(w, string(jj))
}

func check(err error) {

	if err != nil {
		panic( err )
	}

}

func main() {

	// HANDLE JS FILES
	box := rice.MustFindBox("js")
	jsFileServer := http.StripPrefix("/js/", http.FileServer(box.HTTPBox()))
	http.Handle("/js/", jsFileServer)

	boxcss := rice.MustFindBox("css")
	cssFileServer := http.StripPrefix("/css/", http.FileServer(boxcss.HTTPBox()))
	http.Handle("/css/", cssFileServer)

	http.HandleFunc("/operations/", makeHandler(operationsHandler))
	http.HandleFunc("/data/", data_handler)

	http.ListenAndServe(":8080", nil)
}
