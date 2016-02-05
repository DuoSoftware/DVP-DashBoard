// DashBoard project main.go
package main

import (
	"github.com/sukitha/gorest"
	"fmt"
	"net/http"
	"time"
)

func main() {

	//fmt.Println("Hello World!")
	LoadConfiguration()
	InitiateRedis()
	InitiateStatDClient()
	go ClearData()
	gorest.RegisterService(new(DashBoardEvent))
	http.Handle("/", gorest.Handle())
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)

	////fmt.Scanln()
	//client, error := goredis.Dial(&goredis.DialConfig{Address: "127.0.0.1:6379"})
	//if error == nil {
	//client.Set("key", "value", 0, 0, false, false)

	//	go OnEvent("3", "1", "CALLSERVER", "CALL", "IVR", "1111111", client)

	//} else {
	//	fmt.Println(error.Error())
	//}
	//fmt.Println("Hello World!")
}

func ClearData() {
	for {
		tmNow := time.Now().UTC()
		clerTime := time.Date(tmNow.Year(), tmNow.Month(), tmNow.Day(), 23, 59, 59, 0, time.UTC)
		timeToWait := clerTime.Sub(tmNow)
		timer := time.NewTimer(timeToWait)
		<-timer.C
		OnReset()
	}
}
