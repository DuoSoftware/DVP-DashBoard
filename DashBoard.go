package main

import (
	"database/sql"
	"fmt"
	"github.com/fzzy/radix/extra/pubsub"
	"github.com/fzzy/radix/redis"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
	"time"
)

var statClient *StatsdClient

const layout = "2006-01-02T15:04:05Z07:00"

//var log = log4go.NewLogger()

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}

func InitiateStatDClient() {
	host := statsDIp
	port := statsDPort

	//client := statsd.New(host, port)
	statClient = New(host, port)
}

func InitiateRedis() {
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
	//authServer
	authE := c2.Cmd("auth", redisPassword)
	errHndlr(authE.Err)

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

func PersistsMetaData(_class, _type, _category, _window string, count int, _flushEnable, _useSession bool) {
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", pgUser, pgPassword, pgDbname, pgHost, pgPort)
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		fmt.Println(err.Error())
	}

	result, err1 := db.Exec("INSERT INTO \"Dashboard_MetaData\"(\"EventClass\", \"EventType\", \"EventCategory\", \"WindowName\", \"Count\", \"FlushEnable\", \"UseSession\", \"createdAt\", \"updatedAt\") VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", _class, _type, _category, _window, count, _flushEnable, _useSession, time.Now().Local(), time.Now().Local())
	if err1 != nil {
		fmt.Println(err1.Error())
	} else {
		fmt.Println("PersistsMetaData: ", result)
		lInsertedId, err2 := result.LastInsertId()
		fmt.Println(err2)
		fmt.Println("Last inserted Id: ", lInsertedId)
	}
	db.Close()
}

func ReloadMetaData(_class, _type, _category string) bool {
	var result bool
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d", pgUser, pgPassword, pgDbname, pgHost, pgPort)
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		fmt.Println(err.Error())
		result = false
	}

	var EventClass string
	var EventType string
	var EventCategory string
	var WindowName string
	var Count int
	var FlushEnable bool
	var UseSession bool

	err1 := db.QueryRow("SELECT \"EventClass\", \"EventType\", \"EventCategory\", \"WindowName\", \"Count\", \"FlushEnable\", \"UseSession\" FROM \"Dashboard_MetaData\" WHERE \"EventClass\"=$1 AND \"EventType\"=$2 AND \"EventCategory\"=$3", _class, _type, _category).Scan(&EventClass, &EventType, &EventCategory, &WindowName, &Count, &FlushEnable, &UseSession)
	switch {
	case err1 == sql.ErrNoRows:
		fmt.Println("No metaData with that ID.")
		result = false
	case err1 != nil:
		fmt.Println(err1.Error())
		result = false
	default:
		fmt.Printf("EventClass is %s\n", EventClass)
		fmt.Printf("EventType is %s\n", EventType)
		fmt.Printf("EventCategory is %s\n", EventCategory)
		fmt.Printf("WindowName is %s\n", WindowName)
		fmt.Printf("Count is %d\n", Count)
		fmt.Printf("FlushEnable is %t\n", FlushEnable)
		fmt.Printf("UseSession is %t\n", UseSession)
		CacheMetaData(EventClass, EventType, EventCategory, WindowName, Count, FlushEnable, UseSession)
		result = true
	}
	db.Close()
	return result
}

func CacheMetaData(_class, _type, _category, _window string, count int, _flushEnable, _useSession bool) {

	_windowName := fmt.Sprintf("META:%s:%s:%s:WINDOW", _class, _type, _category)
	_incName := fmt.Sprintf("META:%s:%s:%s:COUNT", _class, _type, _category)
	_flushName := fmt.Sprintf("META:%s:%s:%s:FLUSH", _class, _type, _category)
	_useSessionName := fmt.Sprintf("META:%s:%s:%s:USESESSION", _class, _type, _category)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnMeta", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()
	//authServer
	authE := client.Cmd("auth", redisPassword)
	errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	if _flushEnable == true {
		client.Cmd("setnx", _flushName, _window)
	} else {
		client.Cmd("del", _flushName)
	}

	client.Cmd("setnx", _useSessionName, strconv.FormatBool(_useSession))
	client.Cmd("setnx", _windowName, _window)
	client.Cmd("setnx", _incName, strconv.Itoa(count))
}

func OnMeta(_class, _type, _category, _window string, count int, _flushEnable, _useSession bool) {
	CacheMetaData(_class, _type, _category, _window, count, _flushEnable, _useSession)
	PersistsMetaData(_class, _type, _category, _window, count, _flushEnable, _useSession)
}

func OnEvent(_tenent, _company int, _class, _type, _category, _session, _parameter1, _parameter2 string) {
	temp := fmt.Sprintf("Tenant:%d Company:%d Class:%s Type:%s Category:%s Session:%s Param1:%s Param2:%s", _tenent, _company, _class, _type, _category, _session, _parameter1, _parameter2)
	fmt.Println("OnEvent: ", temp)

	_window := fmt.Sprintf("META:%s:%s:%s:WINDOW", _class, _type, _category)
	_inc := fmt.Sprintf("META:%s:%s:%s:COUNT", _class, _type, _category)
	_useSessionName := fmt.Sprintf("META:%s:%s:%s:USESESSION", _class, _type, _category)
	tm := time.Now()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnEvent", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()
	//authServer
	authE := client.Cmd("auth", redisPassword)
	errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	isWindowExist, _ := client.Cmd("exists", _window).Bool()
	isIncExist, _ := client.Cmd("exists", _inc).Bool()

	if isWindowExist == false || isIncExist == false {
		ReloadMetaData(_class, _type, _category)
	}
	window, _werr := client.Cmd("get", _window).Str()
	errHndlr(_werr)
	sinc, _ierr := client.Cmd("get", _inc).Str()
	errHndlr(_ierr)
	useSession, _userr := client.Cmd("get", _useSessionName).Bool()
	errHndlr(_userr)

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
		maxTimeEventName := fmt.Sprintf("MAXTIME:%d:%d:%s,%s,%s", _tenent, _company, window, _parameter1, _parameter2)

		if _parameter1 == "" {
			_parameter1 = "empty"
		}
		countConcStatName := fmt.Sprintf("event.concurrent.%d.%d.%s.%s", _tenent, _company, _parameter1, window)
		gaugeConcStatName := fmt.Sprintf("event.concurrent.%d.%d.%s.%s", _tenent, _company, _parameter1, window)
		timeStatName := fmt.Sprintf("event.timer.%d.%d.%s.%s", _tenent, _company, _parameter1, window)
		totCountStatName := fmt.Sprintf("event.totalcount.%d.%d.%s.%s", _tenent, _company, _parameter1, window)
		totTimeStatName := fmt.Sprintf("event.totaltime.%d.%d.%s.%s", _tenent, _company, _parameter1, window)

		client.Cmd("incr", snapEventName)
		client.Cmd("incr", snapHourlyEventName)

		if iinc > 0 {
			if useSession {
				client.Cmd("hset", sessEventName, "time", tm.Format(layout))
			}
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

				oldMaxTime, _ := client.Cmd("get", maxTimeEventName).Int()
				if oldMaxTime < timeDiff {
					client.Cmd("set", maxTimeEventName, timeDiff)
				}

				statClient.Decrement(countConcStatName)
				statClient.Gauge(gaugeConcStatName, dccount)
				statClient.Gauge(totTimeStatName, rinc)

				duration := int64(tm.UTC().Sub(tm2.UTC()) / time.Millisecond)
				statClient.Timing(timeStatName, duration)
			}

		}

		logWindow := fmt.Sprintf("%s : %s", window, sinc)
		fmt.Println(logWindow)
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
	//authServer
	authE := client.Cmd("auth", redisPassword)
	errHndlr(authE.Err)
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
			maxTimeEventSearch := fmt.Sprintf("MAXTIME:*:%s:*", window)

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

			maxTimeVal, _ := client.Cmd("keys", maxTimeEventSearch).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, maxTimeVal)

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

func OnGetMaxTime(_tenant, _company int, _window, _parameter1, _parameter2 string, resultChannel chan int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetMaxTime", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()
	//authServer
	authE := client.Cmd("auth", redisPassword)
	errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	maxtimeSearch := fmt.Sprintf("MAXTIME:%d:%d:%s,%s,%s", _tenant, _company, _window, _parameter1, _parameter2)
	keyList, _ := client.Cmd("keys", maxtimeSearch).List()
	if len(keyList) > 0 {
		tempMaxTime := 0
		for _, key := range keyList {
			value, _ := client.Cmd("get", key).Int()
			if tempMaxTime < value {
				tempMaxTime = value
			}
		}
		resultChannel <- tempMaxTime

	} else {
		resultChannel <- 0
	}
}

func OnGetCurrentMaxTime(_tenant, _company int, _window, _parameter1, _parameter2 string, resultChannel chan int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetMaxTime", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()
	//authServer
	authE := client.Cmd("auth", redisPassword)
	errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	maxtimeSearch := fmt.Sprintf("SESSION:%d:%d:%s:*:%s:%s", _tenant, _company, _window, _parameter1, _parameter2)
	keyList, _ := client.Cmd("keys", maxtimeSearch).List()
	if len(keyList) > 0 {
		tempMaxTime := 0
		tm := time.Now()
		for _, key := range keyList {
			tmx, _ := client.Cmd("hget", key, "time").Str()
			tm2, _ := time.Parse(layout, tmx)
			timeDiff := int(tm.UTC().Sub(tm2.UTC()).Seconds())
			if tempMaxTime < timeDiff {
				tempMaxTime = timeDiff
			}
		}
		resultChannel <- tempMaxTime

	} else {
		resultChannel <- 0
	}
}

func OnGetCurrentCount(_tenant, _company int, _window, _parameter1, _parameter2 string, resultChannel chan int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetMaxTime", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()
	//authServer
	authE := client.Cmd("auth", redisPassword)
	errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	concurrentSearch := fmt.Sprintf("CONCURRENT:%d:%d:%s:%s:%s", _tenant, _company, _window, _parameter1, _parameter2)
	keyList, _ := client.Cmd("keys", concurrentSearch).List()
	if len(keyList) > 0 {
		temptotal := 0
		for _, key := range keyList {
			value, _ := client.Cmd("get", key).Int()
			temptotal = temptotal + value
		}
		resultChannel <- temptotal

	} else {
		resultChannel <- 0
	}
}

func OnGetAverageTime(_tenant, _company int, _window, _parameter1, _parameter2 string, resultChannel chan float32) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetMaxTime", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()
	//authServer
	authE := client.Cmd("auth", redisPassword)
	errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	totTimeSearch := fmt.Sprintf("TOTALTIME:%d:%d:%s:%s:%s", _tenant, _company, _window, _parameter1, _parameter2)
	totCountSearch := fmt.Sprintf("TOTALCOUNT:%d:%d:%s:%s:%s", _tenant, _company, _window, _parameter1, _parameter2)

	totalTime := 0
	totalCount := 0

	totTimeKeyList, _ := client.Cmd("keys", totTimeSearch).List()
	if len(totTimeKeyList) > 0 {
		temptotal := 0
		for _, key := range totTimeKeyList {
			value, _ := client.Cmd("get", key).Int()
			temptotal = temptotal + value
		}
		totalTime = temptotal

	} else {
		totalTime = 0
	}

	totCountKeyList, _ := client.Cmd("keys", totCountSearch).List()
	if len(totCountKeyList) > 0 {
		temptotal := 0
		for _, key := range totCountKeyList {
			value, _ := client.Cmd("get", key).Int()
			temptotal = temptotal + value
		}
		totalCount = temptotal

	} else {
		totalCount = 0
	}
	fmt.Println("totalTime: ", totalTime)
	fmt.Println("totalCount: ", totalCount)
	avg := float32(totalTime) / float32(totalCount)
	fmt.Println("avg: ", avg)
	resultChannel <- avg
}
