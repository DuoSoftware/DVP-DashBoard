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
	EventClass      string
	EventType       string
	EventCategory   string
	WindowName      string
	Count           int
	FlushEnable     bool
	UseSession      bool
	ThresholdEnable bool
	ThresholdValue  int
}

type Configuration struct {
	RedisIp          string
	RedisPort        string
	RedisDb          int
	ArdsRedisDb      int
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
	ArdsRedisDb      string
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

type QueueDetail struct {
	TotalQueued        int
	TotalAnswered      int
	QueueDropped       int
	CurrentWaiting     int
	MaxWaitTime        int
	CurrentMaxWaitTime int
	AverageWaitTime    float32
}

type QueueDetails struct {
	QueueId   string
	QueueName string
	QueueInfo QueueDetail
}

type SummeryDetail struct {
	Company        int
	Tenant         int
	WindowName     string
	Param1         string
	Param2         string
	MaxTime        int
	TotalCount     int
	TotalTime      int
	ThresholdValue int
	SummaryDate    time.Time
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
	queueDetails       gorest.EndPoint `method:"GET" path:"/QueueDetails/" output:"[]QueueDetails"`
	queueSingleDetail  gorest.EndPoint `method:"GET" path:"/QueueSingleDetail/{queueId:string}" output:"QueueDetails"`
	totalCount         gorest.EndPoint `method:"GET" path:"/TotalCount/{window:string}/{param1:string}/{param2:string}" output:"int"`
}

type DashBoardGraph struct {
	gorest.RestService       `root:"/DashboardGraph/" consumes:"application/json" produces:"application/json"`
	calls                    gorest.EndPoint `method:"GET" path:"/Calls/{duration:int}" output:"string"`
	channels                 gorest.EndPoint `method:"GET" path:"/Channels/{duration:int}" output:"string"`
	bridge                   gorest.EndPoint `method:"GET" path:"/Bridge/{duration:int}" output:"string"`
	queued                   gorest.EndPoint `method:"GET" path:"/Queued/{duration:int}" output:"string"`
	concurrentqueued         gorest.EndPoint `method:"GET" path:"/ConcurrentQueued/{queue:string}/{duration:int}" output:"string"`
	allConcurrentQueued      gorest.EndPoint `method:"GET" path:"/AllConcurrentQueued/{duration:int}" output:"string"`
	newTicket                gorest.EndPoint `method:"GET" path:"/NewTicket/{duration:int}" output:"string"`
	closedTicket             gorest.EndPoint `method:"GET" path:"/ClosedTicket/{duration:int}" output:"string"`
	closedVsOpenTicket       gorest.EndPoint `method:"GET" path:"/ClosedVsOpenTicket/{duration:int}" output:"string"`
	newTicketByUser          gorest.EndPoint `method:"GET" path:"/NewTicketByUser/{duration:int}" output:"string"`
	closedTicketByUser       gorest.EndPoint `method:"GET" path:"/ClosedTicketByUser/{duration:int}" output:"string"`
	closedVsOpenTicketByUser gorest.EndPoint `method:"GET" path:"/ClosedVsOpenTicketByUser/{duration:int}" output:"string"`
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

	go OnMeta(data.EventClass, data.EventType, data.EventCategory, data.WindowName, data.Count, data.FlushEnable, data.UseSession, data.ThresholdEnable, data.ThresholdValue)

	return

}

func (dashboardEvent DashBoardEvent) Reset() {

	const longForm = "Jan 2, 2006 at 3:04pm (MST)"

	go OnReset()

}

func (dashBoardEvent DashBoardEvent) MaxWaiting(window, param1, param2 string) int {
	company, tenant, _ := decodeJwtDashBoardEvent(dashBoardEvent, "dashboardevent", "read")
	if company != 0 && tenant != 0 {
		resultChannel := make(chan int)
		go OnGetMaxTime(tenant, company, window, param1, param2, resultChannel)
		var maxTime = <-resultChannel
		close(resultChannel)
		return maxTime
	} else {
		dashBoardEvent.RB().SetResponseCode(403)
		return 0
	}
}

func (dashBoardEvent DashBoardEvent) CurrentMaxTime(window, param1, param2 string) int {
	company, tenant, _ := decodeJwtDashBoardEvent(dashBoardEvent, "dashboardevent", "read")
	if company != 0 && tenant != 0 {
		resultChannel := make(chan int)
		go OnGetCurrentMaxTime(tenant, company, window, param1, param2, resultChannel)
		var maxTime = <-resultChannel
		close(resultChannel)
		return maxTime
	} else {
		dashBoardEvent.RB().SetResponseCode(403)
		return 0
	}
}

func (dashBoardEvent DashBoardEvent) CurrentCount(window, param1, param2 string) int {
	company, tenant, _ := decodeJwtDashBoardEvent(dashBoardEvent, "dashboardevent", "read")
	if company != 0 && tenant != 0 {
		resultChannel := make(chan int)
		go OnGetCurrentCount(tenant, company, window, param1, param2, resultChannel)
		var maxTime = <-resultChannel
		close(resultChannel)
		return maxTime
	} else {
		dashBoardEvent.RB().SetResponseCode(403)
		return 0
	}
}

func (dashBoardEvent DashBoardEvent) AverageTime(window, param1, param2 string) float32 {
	company, tenant, _ := decodeJwtDashBoardEvent(dashBoardEvent, "dashboardevent", "read")
	if company != 0 && tenant != 0 {
		resultChannel := make(chan float32)
		go OnGetAverageTime(tenant, company, window, param1, param2, resultChannel)
		var maxTime = <-resultChannel
		close(resultChannel)
		fmt.Println("AverageTime: ", maxTime)
		return maxTime
	} else {
		dashBoardEvent.RB().SetResponseCode(403)
		return 0
	}
}

func (dashBoardEvent DashBoardEvent) QueueDetails() []QueueDetails {
	company, tenant, _ := decodeJwtDashBoardEvent(dashBoardEvent, "dashboardevent", "read")
	if company != 0 && tenant != 0 {
		resultChannel := make(chan []QueueDetails)
		go OnGetQueueDetails(tenant, company, resultChannel)
		var queueInfo = <-resultChannel
		close(resultChannel)
		return queueInfo
	} else {
		dashBoardEvent.RB().SetResponseCode(403)
		return make([]QueueDetails, 0)
	}
}

func (dashBoardEvent DashBoardEvent) QueueSingleDetail(queueId string) QueueDetails {
	company, tenant, _ := decodeJwtDashBoardEvent(dashBoardEvent, "dashboardevent", "read")
	if company != 0 && tenant != 0 {
		resultChannel := make(chan QueueDetails)
		go OnGetSingleQueueDetails(tenant, company, queueId, resultChannel)
		var queueInfo = <-resultChannel
		close(resultChannel)
		return queueInfo
	} else {
		dashBoardEvent.RB().SetResponseCode(403)
		var detais = QueueDetails{}
		return detais
	}
}

func (dashBoardEvent DashBoardEvent) TotalCount(window, param1, param2 string) int {
	company, tenant, _ := decodeJwtDashBoardEvent(dashBoardEvent, "dashboardevent", "read")
	fmt.Println(company, tenant)
	if company != 0 && tenant != 0 {
		resultChannel := make(chan int)
		go OnGetTotalCount(tenant, company, window, param1, param2, resultChannel)
		var totalCount = <-resultChannel
		close(resultChannel)
		return totalCount
	} else {
		dashBoardEvent.RB().SetResponseCode(403)
		return 0
	}
}

func (dashBoardGraph DashBoardGraph) Calls(duration int) string {
	company, tenant, _, _ := decodeJwtDashBoardGraph(dashBoardGraph, "dashboardgraph", "read")
	if company != 0 && tenant != 0 {
		resultChannel := make(chan string)
		go OnGetCalls(tenant, company, duration, resultChannel)
		var graphData = <-resultChannel
		close(resultChannel)
		return graphData
	} else {
		dashBoardGraph.RB().SetResponseCode(403)
		return ""
	}
}

func (dashBoardGraph DashBoardGraph) Channels(duration int) string {
	company, tenant, _, _ := decodeJwtDashBoardGraph(dashBoardGraph, "dashboardgraph", "read")
	if company != 0 && tenant != 0 {
		resultChannel := make(chan string)
		go OnGetChannels(tenant, company, duration, resultChannel)
		var graphData = <-resultChannel
		close(resultChannel)
		return graphData
	} else {
		dashBoardGraph.RB().SetResponseCode(403)
		return ""
	}
}

func (dashBoardGraph DashBoardGraph) Bridge(duration int) string {
	company, tenant, _, _ := decodeJwtDashBoardGraph(dashBoardGraph, "dashboardgraph", "read")
	if company != 0 && tenant != 0 {
		resultChannel := make(chan string)
		go OnGetBridge(tenant, company, duration, resultChannel)
		var graphData = <-resultChannel
		close(resultChannel)
		return graphData
	} else {
		dashBoardGraph.RB().SetResponseCode(403)
		return ""
	}
}

func (dashBoardGraph DashBoardGraph) Queued(duration int) string {
	company, tenant, _, _ := decodeJwtDashBoardGraph(dashBoardGraph, "dashboardgraph", "read")
	fmt.Println(company, "::", tenant)
	if company != 0 && tenant != 0 {
		resultChannel := make(chan string)
		go OnGetQueued(tenant, company, duration, resultChannel)
		var graphData = <-resultChannel
		close(resultChannel)
		return graphData
	} else {
		dashBoardGraph.RB().SetResponseCode(403)
		return ""
	}
}

func (dashBoardGraph DashBoardGraph) Concurrentqueued(queue string, duration int) string {
	company, tenant, _, _ := decodeJwtDashBoardGraph(dashBoardGraph, "dashboardgraph", "read")
	if company != 0 && tenant != 0 {
		resultChannel := make(chan string)
		go OnGetConcurrentQueue(tenant, company, duration, queue, resultChannel)
		var graphData = <-resultChannel
		close(resultChannel)
		return graphData
	} else {
		dashBoardGraph.RB().SetResponseCode(403)
		return ""
	}
}

func (dashBoardGraph DashBoardGraph) AllConcurrentQueued(duration int) string {
	company, tenant, _, _ := decodeJwtDashBoardGraph(dashBoardGraph, "dashboardgraph", "read")
	if company != 0 && tenant != 0 {
		resultChannel := make(chan string)
		go OnGetConcurrentQueueTotal(tenant, company, duration, resultChannel)
		var graphData = <-resultChannel
		close(resultChannel)
		return graphData
	} else {
		dashBoardGraph.RB().SetResponseCode(403)
		return ""
	}
}

func (dashBoardGraph DashBoardGraph) NewTicket(duration int) string {
	company, tenant, _, _ := decodeJwtDashBoardGraph(dashBoardGraph, "dashboardgraph", "read")
	fmt.Println(company, "::", tenant)
	if company != 0 && tenant != 0 {
		resultChannel := make(chan string)
		go OnGetTotalNewTicket(tenant, company, duration, resultChannel)
		var graphData = <-resultChannel
		close(resultChannel)
		return graphData
	} else {
		dashBoardGraph.RB().SetResponseCode(403)
		return ""
	}
}

func (dashBoardGraph DashBoardGraph) ClosedTicket(duration int) string {
	company, tenant, _, _ := decodeJwtDashBoardGraph(dashBoardGraph, "dashboardgraph", "read")
	fmt.Println(company, "::", tenant)
	if company != 0 && tenant != 0 {
		resultChannel := make(chan string)
		go OnGetTotalClosedTicket(tenant, company, duration, resultChannel)
		var graphData = <-resultChannel
		close(resultChannel)
		return graphData
	} else {
		dashBoardGraph.RB().SetResponseCode(403)
		return ""
	}
}

func (dashBoardGraph DashBoardGraph) ClosedVsOpenTicket(duration int) string {
	company, tenant, _, _ := decodeJwtDashBoardGraph(dashBoardGraph, "dashboardgraph", "read")
	fmt.Println(company, "::", tenant)
	if company != 0 && tenant != 0 {
		resultChannel := make(chan string)
		go OnGetDiffClosedVsNew(tenant, company, duration, resultChannel)
		var graphData = <-resultChannel
		close(resultChannel)
		return graphData
	} else {
		dashBoardGraph.RB().SetResponseCode(403)
		return ""
	}
}

func (dashBoardGraph DashBoardGraph) NewTicketByUser(duration int) string {
	company, tenant, iss, _ := decodeJwtDashBoardGraph(dashBoardGraph, "dashboardgraph", "read")
	fmt.Println(company, "::", tenant, "::", iss)
	if company != 0 && tenant != 0 {
		resultChannel := make(chan string)
		go OnGetTotalNewTicketByUser(tenant, company, duration, iss, resultChannel)
		var graphData = <-resultChannel
		close(resultChannel)
		return graphData
	} else {
		dashBoardGraph.RB().SetResponseCode(403)
		return ""
	}
}

func (dashBoardGraph DashBoardGraph) ClosedTicketByUser(duration int) string {
	company, tenant, iss, _ := decodeJwtDashBoardGraph(dashBoardGraph, "dashboardgraph", "read")
	fmt.Println(company, "::", tenant, "::", iss)
	if company != 0 && tenant != 0 {
		resultChannel := make(chan string)
		go OnGetTotalClosedTicketByUser(tenant, company, duration, iss, resultChannel)
		var graphData = <-resultChannel
		close(resultChannel)
		return graphData
	} else {
		dashBoardGraph.RB().SetResponseCode(403)
		return ""
	}
}

func (dashBoardGraph DashBoardGraph) ClosedVsOpenTicketByUser(duration int) string {
	company, tenant, iss, _ := decodeJwtDashBoardGraph(dashBoardGraph, "dashboardgraph", "read")
	fmt.Println(company, "::", tenant, "::", iss)
	if company != 0 && tenant != 0 {
		resultChannel := make(chan string)
		go OnGetDiffClosedVsNewByUser(tenant, company, duration, iss, resultChannel)
		var graphData = <-resultChannel
		close(resultChannel)
		return graphData
	} else {
		dashBoardGraph.RB().SetResponseCode(403)
		return ""
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
