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

type PointData map[string][]yamoney.Operation

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


func check_time( time_str string ) time.Time {

	t, err := time.Parse(time.RFC3339, time_str)
	if err != nil {
		log.Println("cannot parse: ", err)
		log.Println("cannot parse: ", time_str)
		return time.Now()
	}
	return t

	//date_regex := regexp.MustCompile(`^(\d{2})/(\d{2})/(\d{4})$`)

	//m := date_regex.FindStringSubmatch(time_str)
	//
	//
	//
	//parsed_month, _ := strconv.ParseInt( m[1], 10, 16 )
	//t := time.Date(m[3], time.Month( parsed_month), m[2], 0, 0, 0, 0, time.Local)
}


func data_handler (w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	from := check_time( r.Form.Get("from") )
	till := check_time( r.Form.Get("till") )

	log.Println("From: ", from)
	log.Println("To: ", till)

	token, err := ioutil.ReadFile(file_name)
	check(err)

	ya, _ := yamoney.NewYaMoney(string(token))

	operations := ya.OperationHistory(
		"from=" + from.Format(time.RFC3339) +
			"&till=" + till.Format(time.RFC3339) +
			"&records=100" )

	var data struct {
		In, Out []interface{}
		Labels map[string]string
		Sum    map[string]float64
		AggregatedOut []interface{}
		Total float64
	}
	var sum float64

	points := make(PointData)

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

		//fmt.Println( op.DateTime, " : ", op.Title, " : ",  op.Amount)
		sum += op.Amount

		curt, _ := time.Parse(time.RFC3339, op.DateTime)

		var rr *[]interface{}
		if op.Direction == "in" {
			rr = &data.In
		} else {
			rr = &data.Out
		}

		kk := curt.Truncate(dur).Unix() * 1000
		k := strconv.FormatInt(curt.Truncate(dur).Unix() * 1000, 10)

		*rr = append( *rr, []interface{}{ kk, op.Amount })

		points[k] = append(points[k], op)

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
		k := strconv.FormatInt(curt.Truncate(dur).Unix() * 1000, 10)


		data.AggregatedOut = append(data.AggregatedOut, []interface{}{ curt.Truncate(dur).Unix() * 1000, data.Sum[k] })
	}

	data.Labels = points.GetLabels()

	log.Println( "Labels: ", data.Labels )

	jj, err := json.Marshal(data)
	check(err)

	log.Println( string(jj) )

	fmt.Fprint(w, string(jj))
}

	func (points PointData) GetLabels() map[string]string {

	timerepl := regexp.MustCompile(".*T")

	var res = make(map[string]string)

	for k := range points {
		ops := points[k]
		//sort.Sort(yamoney.ByDateTime( ops )) // sort !

		// get max lenght of digits placed
		var max_len int = 0
		for o := range ops{
			op := ops[o]

			str := fmt.Sprintf("%0.2f", op.Amount)
			if max_len < len(str) {
				max_len = len(str)
			}
		}

		var sum float64
		for o := range ops {

			op := ops[o]
			//whipe out date
			rdone := timerepl.ReplaceAll([]byte(op.DateTime), []byte(``))

			str := fmt.Sprintf("%6.2f", op.Amount)
			diff := max_len - len(str)

			for i := 1; i <= diff; i++ {
				str = `<span style="color:white">0</span>` + str
			}

			res[k] = fmt.Sprintf(
				//`%s <tr><td><span style="color:green">%s</span></td><td>%f</td></tr>`, data.Labels[k],
				`%s  ★ %s ₽ <span style="color:green">%s </span><span style="color: #78a2b7">%s</span><br/>`, res[k],
				str, string(rdone), op.Title )

			sum += op.Amount
		}

		res[k] = fmt.Sprintf(
			`<br/><span style="background-color:yellow;color:red;font-size:+10">%f</span> <br/>%s`,
			sum, res[k] )
	}

	return res
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
