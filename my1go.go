package main

import (
	"time"
	"log"
	"fmt"
	"flag"
	"./ya-money-go/"
	"io/ioutil"
)

func main() {
	fmt.Println("hello world :)))")

	file_name := flag.String("token_file", "/home/alexander/token_file", "file, containing yamoney token" )
	flag.Parse()

	token, err := ioutil.ReadFile(*file_name)

	if err != nil {
		panic(err)
	}

	ya, _ := yamoney.NewYaMoney(string(token))

	// Start period
	from, _ := time.Parse( "2006-Jan-02 00:00:00", "2016-Aug-01 00:00:00" )
	till := from.AddDate(0, 1, 0)

	//account_info := ya.account_info()
	//account_info := ya.operation_history("type=payment&records=1")
	operations := ya.OperationHistory( "from=" + from.Format(time.RFC3339) + "&till=" + till.Format(time.RFC3339) )

	log.Println(operations)

	var sum float64
	for _, op := range operations {
		fmt.Println(op.Title)
		sum += op.Amount
	}

	log.Println("Sum: ", sum)
}
