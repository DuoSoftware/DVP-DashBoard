// DashBoard project main.go
package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DuoSoftware/gorest"
	"github.com/rs/cors"
)

func main() {

	LoadConfiguration()

	if cacheMachenism == "memory" {

		ReloadAllMetaData()
	}

	InitiateRedis()
	InitiateStatDClient()

	if useMsgQueue == "false" {
		//go PubSub()
		fmt.Println("Redis Pubsub has been removed !!!!")
	} else {
		if useAmqpAdapter == "false" {
			go func() {
				rmqIps := strings.Split(rabbitMQIp, ",")
				currentRmqNodeIndex := 0
				rmqNodeTryCount := 0

				for {
					if len(rmqIps) > 1 {
						if rmqNodeTryCount > 30 {
							fmt.Println("Start to change RMQ node")
							if currentRmqNodeIndex == (len(rmqIps) - 1) {
								currentRmqNodeIndex = 0
							} else {
								currentRmqNodeIndex++
							}
							rmqNodeTryCount = 0
						}
					}
					rmqNodeTryCount++
					fmt.Println("Start Connecting to RMQ: ", rmqIps[currentRmqNodeIndex], " :: TryCount: ", rmqNodeTryCount)
					Worker(rmqIps[currentRmqNodeIndex])
					fmt.Println("End Worker()")
					time.Sleep(2 * time.Second)
				}
			}()
		}
	}

	go StartDecrRetry()
	//go ClearData()

	gorest.RegisterService(new(DashBoardEvent))

	c := cors.New(cors.Options{
		AllowedHeaders: []string{"accept", "authorization"},
	})
	handler := c.Handler(gorest.Handle())
	addr := fmt.Sprintf(":%s", port)
	s := &http.Server{
		Addr:           addr,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.SetKeepAlivesEnabled(false)
	s.ListenAndServe()
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
		OnSetDailyThesholdBreakDown(clerTime)
		OnReset()

		fmt.Println("----------ClearData Wait after reset----------------------")
		timer2 := time.NewTimer(time.Duration(time.Minute * 5))
		<-timer2.C
		fmt.Println("----------End ClearData Wait after reset----------------------")
	}
}

func StartDecrRetry() {
	for {
		ProcessDecrRetry()
		time.Sleep(1 * time.Second)
	}
}
