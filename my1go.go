package main

import (
	"time"
	"log"
	"fmt"
	"flag"
	"./ya-money-go/"
)

func main() {
	fmt.Println("hello world :)))")

	wordPtr := flag.String("word", "foo123", "a string")

	flag.Parse()
	fmt.Println("word:", *wordPtr)

	ya, _ := yamoney.NewYaMoney("410011329603013.1C6223286B618EBF7A24C3EB9EC7D1C6E243789EBBEAD6555D2FB55384C82B0A11614E8F576679C6CB057C9F753716788E7B9F3755753D3BE1BFF1D5369D81FEC1970B5DCBE090BF25F22577CEFDF9F83BA2882763D6AD26E54D41CC31D568B6B5E7E3658AE443805F58AE9BE602C588B7896DF918DE20BFE51638C92555F8AC")

	// Start period
	from, _ := time.Parse( "2006-Jan-02 00:00:00", "2016-Aug-01 00:00:00" )
	till := from.AddDate(0, 1, 0)

	//account_info := ya.account_info()
	//account_info := ya.operation_history("type=payment&records=1")
	operations := ya.OperationHistory( "from=" + from.Format(time.RFC3339) + "&till=" + till.Format(time.RFC3339) )

	log.Println(operations)

	//var res yamoney.ResponseOperations
	//if err := json.Unmarshal([]byte(operations), &res); err != nil {
	//	log.Fatal("FATAL: ", err)
	//}

	//log.Println( "RESPONSE: ", res )

	var sum float64
	for _, op := range operations {
		fmt.Println(op.Title)
		sum += op.Amount
	}

	log.Println("Sum: ", sum)
}
