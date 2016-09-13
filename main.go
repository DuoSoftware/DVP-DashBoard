// DashBoard project main.go
package main

import (
	"fmt"
	//"github.com/DuoSoftware/gorest"
	//"github.com/rs/cors"
	//"net/http"
	"time"
)

func main() {

	//fmt.Println("Hello World!")
	LoadConfiguration()
	//InitiateRedis()
	//InitiateStatDClient()
	ClearData()

	//jwtMiddleware := loadJwtMiddleware()
	//gorest.RegisterService(new(DashBoardEvent))
	//gorest.RegisterService(new(DashBoardGraph))
	//app := jwtMiddleware.Handler(gorest.Handle())
	//c := cors.New(cors.Options{
	//AllowedHeaders: []string{"accept", "authorization"},
	//})
	//handler := c.Handler(gorest.Handle())
	//addr := fmt.Sprintf(":%s", port)
	/*s := &http.Server{
		Addr:           addr,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}*/

	//s.SetKeepAlivesEnabled(false)
	//s.ListenAndServe()
	//http.ListenAndServe(addr, handler)

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
		fmt.Println("----------Start ClearData----------------------")
		location, _ := time.LoadLocation("Asia/Colombo")
		fmt.Println("location:: " + location.String())
		localtime := time.Now().Local()
		fmt.Println("localtime:: " + localtime.String())
		tmNow := time.Now().In(location)
		fmt.Println("tmNow:: " + tmNow.String())
		clerTime := time.Date(tmNow.Year(), tmNow.Month(), tmNow.Day(), 23, 59, 59, 0, location)
		fmt.Println("Next Clear Time:: " + clerTime.String())
		timeToWait := clerTime.Sub(tmNow)
		fmt.Println("timeToWait:: " + timeToWait.String())
		timer := time.NewTimer(timeToWait)
		<-timer.C
		OnSetDailySummary(clerTime)
		OnReset()

		fmt.Println("----------ClearData Wait after reset----------------------")
		timer2 := time.NewTimer(time.Duration(time.Minute * 5))
		<-timer2.C
		fmt.Println("----------End ClearData Wait after reset----------------------")
	}
}
