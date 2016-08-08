//YandexMoneyClient
package yamoney

import (
	"net/http"
	"log"
	"io/ioutil"
	"io"
	"bytes"
	//"github.com/go-telegram-bot-api/telegram-bot-api"
	"encoding/json"
	"time"
)

type Operation struct {
	Title string `json:"title"`
	Direction string `json:"direction"`
	Status string `json:"status"`
	DateTime string `json:"datetime"`
	OperationID string `json:"operation_id"`
	Amount float64 `json:"amount"`
}


type ByAmount []Operation
type ByDateTime []Operation

func (s ByAmount) Len() int {
	return len(s)
}

func (s ByAmount) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByAmount) Less(i, j int) bool {
	return s[i].Amount < s[j].Amount
}

func (s ByDateTime) Len() int {
	return len(s)
}

func (s ByDateTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByDateTime) Less(i, j int) bool {
	return len(s[i].DateTime) < len(s[j].DateTime)
}


type ResponseOperations struct {
	NextRecord string `json:"next_record"`
	Operations []Operation `json:"operations"`
}

type OperationHistoryParams struct {
	from time.Time
	till time.Time
}


type YandexMoneyClient struct {
	token string
	api_url string
}

func ( ya YandexMoneyClient ) account_info() string {
	return ya._execute("account-info", nil)
}

func ( ya YandexMoneyClient ) OperationHistory( params ...string ) []Operation {

	ready_param := ""
	for _, param := range params {
		ready_param += param
	}
	log.Println("READY_PARAM: ", ready_param)

	var next_record string = "0"
	var ops []Operation = []Operation{}
	for i := 0; len(next_record) > 0; i++ {

		operations := ya._execute("operation-history", bytes.NewReader([]byte( ready_param + "&start_record=" + next_record ))  )

		var res ResponseOperations
		if err := json.Unmarshal([]byte(operations), &res); err != nil {
			log.Fatal("FATAL: ", err)
		}

		// TODO: how 2 merge to slices/arrays ?
		for m := range res.Operations {
			ops = append( ops, res.Operations[m] )
		}

		log.Println("NExt REcord: ", res.NextRecord)

		next_record = res.NextRecord
	}

	return ops
}

func NewYaMoney(token string) (YandexMoneyClient, error) {

	ya := YandexMoneyClient{
		token: token,
		api_url: "https://money.yandex.ru/api/",
	}

	return ya, nil
}

func (ya YandexMoneyClient) _execute( cmd string, rbody io.Reader ) string {

	log.Println("Cmd is: ", ya.api_url + cmd)
	//log.Println("BODY: ", rbody)

	req, err := http.NewRequest("POST", ya.api_url + cmd, rbody) // last is post-data
	req.Header.Set("Host", "money.yandex.ru")
	req.Header.Set("Authorization", "Bearer " + ya.token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}
