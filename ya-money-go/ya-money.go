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

	operations := ya._execute("operation-history", bytes.NewReader([]byte( ready_param ))  )

	var res ResponseOperations
	if err := json.Unmarshal([]byte(operations), &res); err != nil {
		log.Fatal("FATAL: ", err)
	}

	return res.Operations
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
