package main

import (
	"fmt"
	"github.com/DuoSoftware/gorest"
	"time"
)

type EventData struct {
	Tenent        int
	Company       int
	EventClass    string
	EventType     string
	EventCategory string
	SessionID     string
	TimeStamp     string
	Parameter1    string
	Parameter2    string
}

type MetaData struct {
	EventClass    string
	EventType     string
	EventCategory string
	WindowName    string
	Count         int
	FlushEnable   bool
	UseSession    bool
}

type Configuration struct {
	RedisIp          string
	RedisPort        string
	RedisDb          int
	RedisPassword    string
	Port             string
	StatsDIp         string
	StatsDPort       int
	PgUser           string
	PgPassword       string
	PgDbname         string
	PgHost           string
	PgPort           int
	SecurityIp       string
	SecurityPort     string
	SecurityPassword string
}

type EnvConfiguration struct {
	RedisIp          string
	RedisPort        string
	RedisDb          string
	RedisPassword    string
	Port             string
	StatsDIp         string
	StatsDPort       string
	PgUser           string
	PgPassword       string
	PgDbname         string
	PgHost           string
	PgPort           string
	SecurityIp       string
	SecurityPort     string
	SecurityPassword string
}

type Result struct {
	Exception     string
	CustomMessage string
	IsSuccess     bool
	Result        string
}

type DashBoardEvent struct {
	gorest.RestService `root:"/DashboardEvent/" consumes:"application/json" produces:"application/json"`
	event              gorest.EndPoint `method:"POST" path:"/Event/" postdata:"EventData"`
	meta               gorest.EndPoint `method:"POST" path:"/Meta/" postdata:"MetaData"`
	reset              gorest.EndPoint `method:"DELETE" path:"/Reset/"`
	maxWaiting         gorest.EndPoint `method:"GET" path:"/MaxWaiting/{window:string}/{param1:string}/{param2:string}" output:"int"`
	currentMaxTime     gorest.EndPoint `method:"GET" path:"/CurrentMaxTime/{window:string}/{param1:string}/{param2:string}" output:"int"`
	currentCount       gorest.EndPoint `method:"GET" path:"/CurrentCount/{window:string}/{param1:string}/{param2:string}" output:"int"`
	averageTime        gorest.EndPoint `method:"GET" path:"/AverageTime/{window:string}/{param1:string}/{param2:string}" output:"float32"`
}

func (dashboardEvent DashBoardEvent) Event(data EventData) {
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	t, _ := time.Parse(longForm, data.TimeStamp)
	fmt.Println(t.Hour(), t.Minute())

	fmt.Println(data.EventClass)
	fmt.Println(data.EventType)
	fmt.Println(data.EventCategory)

	go OnEvent(data.Tenent, data.Company, data.EventClass, data.EventType, data.EventCategory, data.SessionID, data.Parameter1, data.Parameter2)

	return

}

func (dashboardEvent DashBoardEvent) Meta(data MetaData) {

	const longForm = "Jan 2, 2006 at 3:04pm (MST)"

	fmt.Println(data.EventClass)
	fmt.Println(data.EventType)
	fmt.Println(data.EventCategory)

	go OnMeta(data.EventClass, data.EventType, data.EventCategory, data.WindowName, data.Count, data.FlushEnable, data.UseSession)

	return

}

func (dashboardEvent DashBoardEvent) Reset() {

	const longForm = "Jan 2, 2006 at 3:04pm (MST)"

	go OnReset()

}

func (dashBoardEvent DashBoardEvent) MaxWaiting(window, param1, param2 string) int {
	company, tenant := validateCompanyTenant(dashBoardEvent)
	if company != 0 && tenant != 0 {
		resultChannel := make(chan int)
		go OnGetMaxTime(tenant, company, window, param1, param2, resultChannel)
		var maxTime = <-resultChannel
		close(resultChannel)
		return maxTime
	} else {
		return 0
	}
}

func (dashBoardEvent DashBoardEvent) CurrentMaxTime(window, param1, param2 string) int {
	company, tenant := validateCompanyTenant(dashBoardEvent)
	if company != 0 && tenant != 0 {
		resultChannel := make(chan int)
		go OnGetCurrentMaxTime(tenant, company, window, param1, param2, resultChannel)
		var maxTime = <-resultChannel
		close(resultChannel)
		return maxTime
	} else {
		return 0
	}
}

func (dashBoardEvent DashBoardEvent) CurrentCount(window, param1, param2 string) int {
	company, tenant := validateCompanyTenant(dashBoardEvent)
	if company != 0 && tenant != 0 {
		resultChannel := make(chan int)
		go OnGetCurrentCount(tenant, company, window, param1, param2, resultChannel)
		var maxTime = <-resultChannel
		close(resultChannel)
		return maxTime
	} else {
		return 0
	}
}

func (dashBoardEvent DashBoardEvent) AverageTime(window, param1, param2 string) float32 {
	company, tenant := validateCompanyTenant(dashBoardEvent)
	if company != 0 && tenant != 0 {
		resultChannel := make(chan float32)
		go OnGetAverageTime(tenant, company, window, param1, param2, resultChannel)
		var maxTime = <-resultChannel
		close(resultChannel)
		return maxTime
	} else {
		return 0
	}
}

//Registerres a service on the rootpath.
//See example below:
//
//	package main
//	import (
// 	   "code.google.com/p/gorest"
//	        "http"
//	)
//	func main() {
//	    gorest.RegisterService(new(HelloService)) //Register our service
//	    http.Handle("/",gorest.Handle())
// 	   http.ListenAndServe(":8787",nil)
//	}
//
//	//Service Definition
//	type HelloService struct {
//	    gorest.RestService `root:"/tutorial/"`
//	    helloWorld  gorest.EndPoint `method:"GET" path:"/hello-world/" output:"string"`
//	    sayHello    gorest.EndPoint `method:"GET" path:"/hello/{name:string}" output:"string"`
//	}
//	func(serv HelloService) HelloWorld() string{
// 	   return "Hello World"
//	}
//	func(serv HelloService) SayHello(name string) string{
//	    return "Hello " + name
//	}
