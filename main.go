// DashBoard project main.go
package main

import (
	"fmt"
	"github.com/DuoSoftware/gorest"
	"github.com/rs/cors"
	"net/http"
	"time"
)

func main() {

	//fmt.Println("Hello World!")
	LoadConfiguration()
	InitiateRedis()
	InitiateStatDClient()
	go ClearData()

	jwtMiddleware := loadJwtMiddleware()
	gorest.RegisterService(new(DashBoardEvent))
	gorest.RegisterService(new(DashBoardGraph))
	app := jwtMiddleware.Handler(gorest.Handle())
	c := cors.New(cors.Options{
		AllowedHeaders: []string{"accept", "authorization"},
	})
	handler := c.Handler(app)
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, handler)

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
		OnSetDailySummary(clerTime)
		OnReset()
	}
}
