package main

import (
	//"code.google.com/p/log4go"
	"encoding/json"
	"fmt"
	"github.com/fzzy/radix/extra/pubsub"
	"github.com/fzzy/radix/redis"
	"os"
	"strconv"
	"strings"
	"time"
)

var redisIp string
var redisDb int
var writeConfig bool
var statClient *StatsdClient

const layout = "2006-01-02T15:04:05Z07:00"

//var log = log4go.NewLogger()

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}

func InitiateStatDClient() {
	host := "45.55.142.207"
	port := 8125

	//client := statsd.New(host, port)
	statClient = New(host, port)
}

func InitiateRedis() {
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
		redisIp = "127.0.0.1:6379"
		redisDb = 8
		writeConfig = false
	} else {
		redisIp = configuration.RedisIp
		redisDb = configuration.RedisDb
		writeConfig = configuration.WriteLogFile
	}

	fmt.Println("RedisIp:", redisIp)
	fmt.Println("RedisDb:", redisDb)
	go PubSub()
}

func PubSub() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in PubSub", r)
		}
	}()
	c2, err := redis.Dial("tcp", redisIp)
	errHndlr(err)
	defer c2.Close()
	psc := pubsub.NewSubClient(c2)
	psr := psc.Subscribe("events")
	ppsr := psc.PSubscribe("EVENT:*")

	if ppsr.Err == nil {

		for {
			psr = psc.Receive()
			if psr.Err != nil {

				fmt.Println(psr.Err.Error())

				break
			}
			list := strings.Split(psr.Message, ":")
			fmt.Println(list)
			if len(list) >= 8 {
				stenent := list[1]
				scompany := list[2]
				sclass := list[3]
				stype := list[4]
				scategory := list[5]
				sparam1 := list[6]
				sparam2 := list[7]
				ssession := list[8]

				itenet, _ := strconv.Atoi(stenent)
				icompany, _ := strconv.Atoi(scompany)

				OnEvent(itenet, icompany, sclass, stype, scategory, ssession, sparam1, sparam2)
			}

		}
		//s := strings.Split("127.0.0.1:5432", ":")
	}

	psc.Unsubscribe("events")

}

func OnMeta(_class, _type, _category, _window string, count int, _flushEnable bool) {

	_windowName := fmt.Sprintf("META:%s:%s:%s:WINDOW", _class, _type, _category)
	_incName := fmt.Sprintf("META:%s:%s:%s:COUNT", _class, _type, _category)
	_flushName := fmt.Sprintf("META:%s:%s:%s:FLUSH", _class, _type, _category)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnMeta", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	if _flushEnable == true {
		client.Cmd("setnx", _flushName, _window)
	} else {
		client.Cmd("del", _flushName)
	}

	client.Cmd("setnx", _windowName, _window)
	client.Cmd("setnx", _incName, strconv.Itoa(count))
}

func OnEvent(_tenent, _company int, _class, _type, _category, _session, _parameter1, _parameter2 string) {
	temp := fmt.Sprintf("Tenant:%d Company:%d Class:%s Type:%s Category:%s Session:%s Param1:%s Param2:%s", _tenent, _company, _class, _type, _category, _session, _parameter1, _parameter2)
	fmt.Println("OnEvent: ", temp)

	_window := fmt.Sprintf("META:%s:%s:%s:WINDOW", _class, _type, _category)
	_inc := fmt.Sprintf("META:%s:%s:%s:COUNT", _class, _type, _category)
	tm := time.Now()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnEvent", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	window, _werr := client.Cmd("get", _window).Str()
	errHndlr(_werr)
	sinc, _ierr := client.Cmd("get", _inc).Str()
	errHndlr(_ierr)

	iinc, berr := strconv.Atoi(sinc)

	fmt.Println("iinc value is ", iinc)

	if _werr == nil && _ierr == nil && berr == nil {
		//if writeConfig {
		//	log.AddFilter("file", log4go.FINE, log4go.NewFileLogWriter("C://IISLOGS/DashBoardGo.txt", false))
		//	logData := fmt.Sprintf("Tenant:%d Company:%d Class:%s Type:%s Category:%s Session:%s Param1:%s Param2:%s Window:%s", _tenent, _company, _class, _type, _category, _session, _parameter1, _parameter2, window)

		//	log.Info("------------------------------------------\n")
		//	log.Info(logData, "\n")
		//	log.Info("iinc value is ", fmt.Sprintf("%d", iinc), "\n")
		//	log.Close()
		//}

		snapEventName := fmt.Sprintf("SNAPSHOT:%d:%d:%s:%s:%s:%s:%d:%d", _tenent, _company, window, _class, _type, _category, tm.Hour(), tm.Minute())
		snapHourlyEventName := fmt.Sprintf("SNAPSHOTHOURLY:%d:%d:%s:%s:%s:%s:%d", _tenent, _company, window, _class, _type, _category, tm.Hour())
		concEventName := fmt.Sprintf("CONCURRENT:%d:%d:%s:%s:%s", _tenent, _company, window, _parameter1, _parameter2)
		sessEventName := fmt.Sprintf("SESSION:%d:%d:%s:%s:%s:%s", _tenent, _company, window, _session, _parameter1, _parameter2)
		totTimeEventName := fmt.Sprintf("TOTALTIME:%d:%d:%s:%s:%s", _tenent, _company, window, _parameter1, _parameter2)
		totCountEventName := fmt.Sprintf("TOTALCOUNT:%d:%d:%s:%s:%s", _tenent, _company, window, _parameter1, _parameter2)
		totCountHrEventName := fmt.Sprintf("TOTALCOUNTHR:%d:%d:%s:%s:%s:%d:%d", _tenent, _company, window, _parameter1, _parameter2, tm.Hour(), tm.Minute())

		countConcStatName := fmt.Sprintf("event.concurrent.%d.%d.%s", _tenent, _company, window)
		gaugeConcStatName := fmt.Sprintf("event.concurrent.%d.%d.%s", _tenent, _company, window)
		timeStatName := fmt.Sprintf("event.timer.%d.%d.%s", _tenent, _company, window)
		totCountStatName := fmt.Sprintf("event.totalcount.%d.%d.%s", _tenent, _company, window)
		totTimeStatName := fmt.Sprintf("event.totaltime.%d.%d.%s", _tenent, _company, window)

		client.Cmd("incr", snapEventName)
		client.Cmd("incr", snapHourlyEventName)

		if iinc > 0 {
			client.Cmd("hset", sessEventName, "time", tm.Format(layout))
			ccount, _ := client.Cmd("incr", concEventName).Int()
			tcount, _ := client.Cmd("incr", totCountEventName).Int()
			client.Cmd("incr", totCountHrEventName)

			fmt.Println("tcount ", tcount)
			fmt.Println("ccount ", ccount)

			statClient.Increment(countConcStatName)
			statClient.Gauge(gaugeConcStatName, ccount)
			statClient.Gauge(totCountStatName, tcount)
			fmt.Println("tcount ", tcount)

		} else {
			tmx, _ := client.Cmd("hget", sessEventName, "time").Str()

			tm2, _ := time.Parse(layout, tmx)

			timeDiff := int(tm.UTC().Sub(tm2.UTC()).Seconds())

			fmt.Println(timeDiff)

			isdel, _ := client.Cmd("del", sessEventName).Int()
			if isdel == 1 {
				rinc, _ := client.Cmd("incrby", totTimeEventName, timeDiff).Int()
				dccount, _ := client.Cmd("decr", concEventName).Int()

				statClient.Decrement(countConcStatName)
				statClient.Gauge(gaugeConcStatName, dccount)
				statClient.Gauge(totTimeStatName, rinc)

				duration := int64(tm.UTC().Sub(tm2.UTC()) / time.Millisecond)
				statClient.Timing(timeStatName, duration)
			}

		}

		fmt.Sprintf("%s %s", window, sinc)

	}

}

func OnReset() {

	_searchName := fmt.Sprintf("META:*:FLUSH")
	fmt.Println("Search Windows to Flush: ", _searchName)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnReset", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	val, _ := client.Cmd("keys", _searchName).List()
	_windowList := make([]string, 0)
	_keysToRemove := make([]string, 0)
	lenth := len(val)
	fmt.Println(lenth)
	if lenth > 0 {
		for _, value := range val {
			tmx, _ := client.Cmd("get", value).Str()

			_windowList = AppendIfMissing(_windowList, tmx)
		}

		for _, window := range _windowList {

			fmt.Println("WindowList_: ", window)

			snapEventSearch := fmt.Sprintf("SNAPSHOT:*:%s:*", window)
			snapHourlyEventSearch := fmt.Sprintf("SNAPSHOTHOURLY:*:%s:*", window)
			concEventSearch := fmt.Sprintf("CONCURRENT:*:%s:*", window)
			sessEventSearch := fmt.Sprintf("SESSION:*:%s:*", window)
			totTimeEventSearch := fmt.Sprintf("TOTALTIME:*:%s:*", window)
			totCountEventSearch := fmt.Sprintf("TOTALCOUNT:*:%s:*", window)
			totCountHr := fmt.Sprintf("TOTALCOUNTHR:*:%s:*", window)

			snapVal, _ := client.Cmd("keys", snapEventSearch).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, snapVal)

			snapHourlyVal, _ := client.Cmd("keys", snapHourlyEventSearch).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, snapHourlyVal)

			concVal, _ := client.Cmd("keys", concEventSearch).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, concVal)

			sessVal, _ := client.Cmd("keys", sessEventSearch).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, sessVal)

			totTimeVal, _ := client.Cmd("keys", totTimeEventSearch).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, totTimeVal)

			totCountVal, _ := client.Cmd("keys", totCountEventSearch).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, totCountVal)

			totCountHrVal, _ := client.Cmd("keys", totCountHr).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, totCountHrVal)

		}

		for _, remove := range _keysToRemove {
			fmt.Println("remove_: ", remove)
			client.Cmd("del", remove)
		}
	}
}

func AppendIfMissing(windowList []string, i string) []string {
	for _, ele := range windowList {
		if ele == i {
			return windowList
		}
	}
	return append(windowList, i)
}

func AppendListIfMissing(windowList1 []string, windowList2 []string) []string {
	notExist := true
	for _, ele2 := range windowList2 {
		for _, ele := range windowList1 {
			if ele == ele2 {
				notExist = false
				break
			}
		}

		if notExist == true {
			windowList1 = append(windowList1, ele2)
		}
	}

	return windowList1
}
