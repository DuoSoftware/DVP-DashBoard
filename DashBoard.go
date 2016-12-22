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

var dashboardMetaInfo []MetaData

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
	c2, err := redis.Dial("tcp", redisPubSubIp)
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

				go OnEvent(itenet, icompany, sclass, stype, scategory, ssession, sparam1, sparam2)
			}

		}
		//s := strings.Split("127.0.0.1:5432", ":")
	}

	psc.Unsubscribe("events")

}

func PersistsSummaryData(_summary SummeryDetail) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in PersistsSummaryData", r)
		}
	}()
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", pgUser, pgPassword, pgDbname, pgHost, pgPort)
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		fmt.Println(err.Error())
	}

	result, err1 := db.Exec("INSERT INTO \"Dashboard_DailySummaries\"(\"Company\", \"Tenant\", \"WindowName\", \"Param1\", \"Param2\", \"MaxTime\", \"TotalCount\", \"TotalTime\", \"ThresholdValue\", \"SummaryDate\", \"createdAt\", \"updatedAt\") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", _summary.Company, _summary.Tenant, _summary.WindowName, _summary.Param1, _summary.Param2, _summary.MaxTime, _summary.TotalCount, _summary.TotalTime, _summary.ThresholdValue, _summary.SummaryDate, time.Now().Local(), time.Now().Local())
	if err1 != nil {
		fmt.Println(err1.Error())
	} else {
		fmt.Println("PersistsSummaryData: ", result)
		lInsertedId, err2 := result.LastInsertId()
		fmt.Println(err2)
		fmt.Println("Last inserted Id: ", lInsertedId)
	}
	db.Close()
}

func PersistsThresholdBreakDown(_summary ThresholdBreakDownDetail) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in PersistsThresholdBreakDown", r)
		}
	}()
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", pgUser, pgPassword, pgDbname, pgHost, pgPort)
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		fmt.Println(err.Error())
	}

	result, err1 := db.Exec("INSERT INTO \"Dashboard_ThresholdBreakDowns\"(\"Company\", \"Tenant\", \"WindowName\", \"Param1\", \"Param2\", \"BreakDown\", \"ThresholdCount\", \"SummaryDate\", \"Hour\", \"createdAt\", \"updatedAt\") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", _summary.Company, _summary.Tenant, _summary.WindowName, _summary.Param1, _summary.Param2, _summary.BreakDown, _summary.ThresholdCount, _summary.SummaryDate, _summary.Hour, time.Now().Local(), time.Now().Local())
	if err1 != nil {
		fmt.Println(err1.Error())
	} else {
		fmt.Println("PersistsThresholdBreakDown: ", result)
		lInsertedId, err2 := result.LastInsertId()
		fmt.Println(err2)
		fmt.Println("Last inserted Id: ", lInsertedId)
	}
	db.Close()
}

func PersistsMetaData(_class, _type, _category, _window string, count int, _flushEnable, _useSession, _thresholdEnable bool, _thresholdValue int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in PersistsMetaData", r)
		}
	}()
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", pgUser, pgPassword, pgDbname, pgHost, pgPort)
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		fmt.Println(err.Error())
	}

	result, err1 := db.Exec("INSERT INTO \"Dashboard_MetaData\"(\"EventClass\", \"EventType\", \"EventCategory\", \"WindowName\", \"Count\", \"FlushEnable\", \"UseSession\", \"ThresholdEnable\", \"ThresholdValue\", \"createdAt\", \"updatedAt\") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", _class, _type, _category, _window, count, _flushEnable, _useSession, _thresholdEnable, _thresholdValue, time.Now().Local(), time.Now().Local())
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
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in ReloadMetaData", r)
		}
	}()
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
	var ThresholdEnable bool
	var ThresholdValue int

	err1 := db.QueryRow("SELECT \"EventClass\", \"EventType\", \"EventCategory\", \"WindowName\", \"Count\", \"FlushEnable\", \"UseSession\", \"ThresholdEnable\", \"ThresholdValue\" FROM \"Dashboard_MetaData\" WHERE \"EventClass\"=$1 AND \"EventType\"=$2 AND \"EventCategory\"=$3", _class, _type, _category).Scan(&EventClass, &EventType, &EventCategory, &WindowName, &Count, &FlushEnable, &UseSession, &ThresholdEnable, &ThresholdValue)
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
		fmt.Printf("ThresholdEnable is %t\n", ThresholdEnable)
		fmt.Printf("ThresholdValue is %d\n", ThresholdValue)
		CacheMetaData(EventClass, EventType, EventCategory, WindowName, Count, FlushEnable, UseSession, ThresholdEnable, ThresholdValue)
		result = true
	}
	db.Close()
	return result
}

func ReloadAllMetaData() bool {
	fmt.Println("-------------------------ReloadAllMetaData----------------------")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in ReloadAllMetaData", r)
		}
	}()
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
	var ThresholdEnable bool
	var ThresholdValue int

	//err1 := db.QueryRow("SELECT \"EventClass\", \"EventType\", \"EventCategory\", \"WindowName\", \"Count\", \"FlushEnable\", \"UseSession\", \"ThresholdEnable\", \"ThresholdValue\" FROM \"Dashboard_MetaData\"").Scan(&EventClass, &EventType, &EventCategory, &WindowName, &Count, &FlushEnable, &UseSession, &ThresholdEnable, &ThresholdValue)
	dataRows, err1 := db.Query("SELECT \"EventClass\", \"EventType\", \"EventCategory\", \"WindowName\", \"Count\", \"FlushEnable\", \"UseSession\", \"ThresholdEnable\", \"ThresholdValue\" FROM \"Dashboard_MetaData\"")
	switch {
	case err1 == sql.ErrNoRows:
		fmt.Println("No metaData with that ID.")
		result = false
	case err1 != nil:
		fmt.Println(err1.Error())
		result = false
	default:
		dashboardMetaInfo = make([]MetaData, 0)
		for dataRows.Next() {
			dataRows.Scan(&EventClass, &EventType, &EventCategory, &WindowName, &Count, &FlushEnable, &UseSession, &ThresholdEnable, &ThresholdValue)

			fmt.Printf("EventClass is %s\n", EventClass)
			fmt.Printf("EventType is %s\n", EventType)
			fmt.Printf("EventCategory is %s\n", EventCategory)
			fmt.Printf("WindowName is %s\n", WindowName)
			fmt.Printf("Count is %d\n", Count)
			fmt.Printf("FlushEnable is %t\n", FlushEnable)
			fmt.Printf("UseSession is %t\n", UseSession)
			fmt.Printf("ThresholdEnable is %t\n", ThresholdEnable)
			fmt.Printf("ThresholdValue is %d\n", ThresholdValue)

			if cacheMachenism == "redis" {
				CacheMetaData(EventClass, EventType, EventCategory, WindowName, Count, FlushEnable, UseSession, ThresholdEnable, ThresholdValue)
			} else {
				var mData MetaData
				mData.EventClass = EventClass
				mData.EventType = EventType
				mData.EventCategory = EventCategory
				mData.Count = Count
				mData.FlushEnable = FlushEnable
				mData.ThresholdEnable = ThresholdEnable
				mData.ThresholdValue = ThresholdValue
				mData.UseSession = UseSession
				mData.WindowName = WindowName

				dashboardMetaInfo = append(dashboardMetaInfo, mData)
			}
		}
		dataRows.Close()
		result = true
	}
	db.Close()
	fmt.Println("DashBoard MetaData:: ", dashboardMetaInfo)
	return result
}

func CacheMetaData(_class, _type, _category, _window string, count int, _flushEnable, _useSession, _thresholdEnable bool, _thresholdValue int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in CacheMetaData", r)
		}
	}()
	_windowName := fmt.Sprintf("META:%s:%s:%s:WINDOW", _class, _type, _category)
	_incName := fmt.Sprintf("META:%s:%s:%s:COUNT", _class, _type, _category)
	_flushName := fmt.Sprintf("META:%s:%s:%s:FLUSH", _class, _type, _category)
	_useSessionName := fmt.Sprintf("META:%s:%s:%s:USESESSION", _class, _type, _category)
	_thresholdEnableName := fmt.Sprintf("META:%s:%s:%s:thresholdEnable", _class, _type, _category)

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

	if _thresholdEnable == true {
		client.Cmd("setnx", _thresholdEnableName, _thresholdValue)
	} else {
		client.Cmd("del", _thresholdEnableName)
	}

	client.Cmd("setnx", _useSessionName, strconv.FormatBool(_useSession))
	client.Cmd("setnx", _windowName, _window)
	client.Cmd("setnx", _incName, strconv.Itoa(count))
}

func OnMeta(_class, _type, _category, _window string, count int, _flushEnable, _useSession, _thresholdEnable bool, _thresholdValue int) {
	CacheMetaData(_class, _type, _category, _window, count, _flushEnable, _useSession, _thresholdEnable, _thresholdValue)
	PersistsMetaData(_class, _type, _category, _window, count, _flushEnable, _useSession, _thresholdEnable, _thresholdValue)
}

func OnEvent(_tenent, _company int, _class, _type, _category, _session, _parameter1, _parameter2 string) {

	temp := fmt.Sprintf("Tenant:%d Company:%d Class:%s Type:%s Category:%s Session:%s Param1:%s Param2:%s", _tenent, _company, _class, _type, _category, _session, _parameter1, _parameter2)
	fmt.Println("OnEvent: ", temp)

	location, _ := time.LoadLocation("Asia/Colombo")
	fmt.Println("location:: " + location.String())

	tm := time.Now().In(location)
	fmt.Println("tmNow:: " + tm.String())

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

	var window, sinc, useSession, threshold string
	var iinc int
	var thresholdEnabled bool
	var _werr, _ierr, _userr, _thresherr, berr error

	if cacheMachenism == "redis" {
		fmt.Println("---------------------Use Redis----------------------")

		_window := fmt.Sprintf("META:%s:%s:%s:WINDOW", _class, _type, _category)
		_inc := fmt.Sprintf("META:%s:%s:%s:COUNT", _class, _type, _category)
		_useSessionName := fmt.Sprintf("META:%s:%s:%s:USESESSION", _class, _type, _category)
		_thresholdEnableName := fmt.Sprintf("META:%s:%s:%s:thresholdEnable", _class, _type, _category)

		isWindowExist, _ := client.Cmd("exists", _window).Bool()
		isIncExist, _ := client.Cmd("exists", _inc).Bool()

		if isWindowExist == false || isIncExist == false {
			ReloadMetaData(_class, _type, _category)
		}
		window, _werr = client.Cmd("get", _window).Str()
		errHndlr(_werr)
		sinc, _ierr = client.Cmd("get", _inc).Str()
		errHndlr(_ierr)
		useSession, _userr = client.Cmd("get", _useSessionName).Str()
		errHndlr(_userr)
		threshold, _thresherr = client.Cmd("get", _thresholdEnableName).Str()
		errHndlr(_thresherr)

		if threshold != "" {
			thresholdEnabled = true
		} else {
			thresholdEnabled = false
		}

		iinc, berr = strconv.Atoi(sinc)

	} else {
		fmt.Println("---------------------Use Memoey----------------------")
		for _, dmi := range dashboardMetaInfo {
			if dmi.EventClass == _class && dmi.EventType == _type && dmi.EventCategory == _category {
				window = dmi.WindowName
				iinc = dmi.Count
				useSession = strconv.FormatBool(dmi.UseSession)
				threshold = strconv.Itoa(dmi.ThresholdValue)
				thresholdEnabled = dmi.ThresholdEnable
				break
			}
		}
	}

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
		maxTimeEventName := fmt.Sprintf("MAXTIME:%d:%d:%s:%s:%s", _tenent, _company, window, _parameter1, _parameter2)
		thresholdEventName := fmt.Sprintf("THRESHOLD:%d:%d:%s:%s:%s", _tenent, _company, window, _parameter1, _parameter2)
		thresholdBreakDownEventName := fmt.Sprintf("THRESHOLDBREAKDOWN:%d:%d:%s:%s:%s", _tenent, _company, window, _parameter1, _parameter2)

		if _parameter1 == "" {
			_parameter1 = "empty"
		}

		var statsDPath string
		switch _class {
		case "TICKET":
			statsDPath = "ticket"
			break

		default:
			statsDPath = "common"
			break
		}

		countConcStatName := fmt.Sprintf("event.%s.concurrent.%d.%d.%s.%s", statsDPath, _tenent, _company, _parameter1, window)
		gaugeConcStatName := fmt.Sprintf("event.%s.concurrent.%d.%d.%s.%s", statsDPath, _tenent, _company, _parameter1, window)
		timeStatName := fmt.Sprintf("event.%s.timer.%d.%d.%s.%s", statsDPath, _tenent, _company, _parameter1, window)
		totCountStatName := fmt.Sprintf("event.%s.totalcount.%d.%d.%s.%s", statsDPath, _tenent, _company, _parameter1, window)
		totTimeStatName := fmt.Sprintf("event.%s.totaltime.%d.%d.%s.%s", statsDPath, _tenent, _company, _parameter1, window)

		client.Cmd("incr", snapEventName)
		client.Cmd("incr", snapHourlyEventName)

		if iinc > 0 {
			if useSession == "true" {
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
			sessEventSearch := fmt.Sprintf("SESSION:%d:%d:%s:%s:*", _tenent, _company, window, _session)
			sessEvents, _ := client.Cmd("keys", sessEventSearch).List()
			if len(sessEvents) > 0 {
				tmx, _ := client.Cmd("hget", sessEvents[0], "time").Str()
				tm2, _ := time.Parse(layout, tmx)
				timeDiff := int(tm.Sub(tm2.In(location)).Seconds())

				if timeDiff < 0 {
					timeDiff = 0
				}

				fmt.Println(timeDiff)

				isdel, _ := client.Cmd("del", sessEvents[0]).Int()
				if isdel == 1 {
					rinc, _ := client.Cmd("incrby", totTimeEventName, timeDiff).Int()
					dccount, _ := client.Cmd("decr", concEventName).Int()

					if dccount < 0 {
						fmt.Println("reset minus concurrent count:: incr by 1 :: ", concEventName)
						dccount, _ = client.Cmd("incr", concEventName).Int()
					}

					oldMaxTime, _ := client.Cmd("get", maxTimeEventName).Int()
					if oldMaxTime < timeDiff {
						client.Cmd("set", maxTimeEventName, timeDiff)
					}
					if window != "QUEUE" {
						statClient.Decrement(countConcStatName)
					}
					if thresholdEnabled == true && threshold != "" {
						thValue, _ := strconv.Atoi(threshold)

						if thValue > 0 {
							thHour := tm.Hour()

							if timeDiff > thValue {
								thcount, _ := client.Cmd("incr", thresholdEventName).Int()
								fmt.Println(thresholdEventName, ": ", thcount)

								thValue_2 := thValue * 2
								thValue_4 := thValue * 4
								thValue_8 := thValue * 8
								thValue_10 := thValue * 10
								thValue_12 := thValue * 12

								fmt.Println("thValue_2::", thValue_2)
								fmt.Println("thValue_4::", thValue_4)
								fmt.Println("thValue_8::", thValue_8)
								fmt.Println("thValue_10::", thValue_10)
								fmt.Println("thValue_12::", thValue_12)

								if timeDiff > thValue && timeDiff <= thValue_2 {
									thresholdBreakDown_1 := fmt.Sprintf("%s:%d:%d:%d", thresholdBreakDownEventName, thHour, thValue, thValue_2)
									client.Cmd("incr", thresholdBreakDown_1)
									fmt.Println("thresholdBreakDown_1::", thresholdBreakDown_1)
								} else if timeDiff > thValue_2 && timeDiff <= thValue_4 {
									thresholdBreakDown_2 := fmt.Sprintf("%s:%d:%d:%d", thresholdBreakDownEventName, thHour, thValue_2, thValue_4)
									client.Cmd("incr", thresholdBreakDown_2)
									fmt.Println("thresholdBreakDown_2::", thresholdBreakDown_2)
								} else if timeDiff > thValue_4 && timeDiff <= thValue_8 {
									thresholdBreakDown_3 := fmt.Sprintf("%s:%d:%d:%d", thresholdBreakDownEventName, thHour, thValue_4, thValue_8)
									client.Cmd("incr", thresholdBreakDown_3)
									fmt.Println("thresholdBreakDown_3::", thresholdBreakDown_3)
								} else if timeDiff > thValue_8 && timeDiff <= thValue_10 {
									thresholdBreakDown_4 := fmt.Sprintf("%s:%d:%d:%d", thresholdBreakDownEventName, thHour, thValue_8, thValue_10)
									client.Cmd("incr", thresholdBreakDown_4)
									fmt.Println("thresholdBreakDown_4::", thresholdBreakDown_4)
								} else if timeDiff > thValue_10 && timeDiff <= thValue_12 {
									thresholdBreakDown_5 := fmt.Sprintf("%s:%d:%d:%d", thresholdBreakDownEventName, thHour, thValue_10, thValue_12)
									client.Cmd("incr", thresholdBreakDown_5)
									fmt.Println("thresholdBreakDown_5::", thresholdBreakDown_5)
								} else {
									thresholdBreakDown_6 := fmt.Sprintf("%s:%d:%d:%s", thresholdBreakDownEventName, thHour, thValue_12, "gt")
									client.Cmd("incr", thresholdBreakDown_6)
									fmt.Println("thresholdBreakDown_6::", thresholdBreakDown_6)
								}
							} else {
								thresholdBreakDown_7 := fmt.Sprintf("%s:%d:%s:%d", thresholdBreakDownEventName, thHour, "lt", thValue)
								client.Cmd("incr", thresholdBreakDown_7)
								fmt.Println("thresholdBreakDown_7::", thresholdBreakDown_7)
							}
						}
					}
					statClient.Gauge(gaugeConcStatName, dccount)
					statClient.Gauge(totTimeStatName, rinc)

					duration := int64(tm.Sub(tm2.In(location)) / time.Millisecond)
					statClient.Timing(timeStatName, duration)
				}
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
	_loginSessions := make([]string, 0)
	_productivitySessions := make([]string, 0)
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
			thresholdEventSearch := fmt.Sprintf("THRESHOLD:*:%s:*", window)
			thresholdBDEventSearch := fmt.Sprintf("THRESHOLDBREAKDOWN:*:%s:*", window)

			snapVal, _ := client.Cmd("keys", snapEventSearch).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, snapVal)

			snapHourlyVal, _ := client.Cmd("keys", snapHourlyEventSearch).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, snapHourlyVal)

			concVal, _ := client.Cmd("keys", concEventSearch).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, concVal)

			sessVal, _ := client.Cmd("keys", sessEventSearch).List()
			for _, sess := range sessVal {
				sessItems := strings.Split(sess, ":")
				if len(sessItems) >= 4 && sessItems[3] == "LOGIN" {
					_loginSessions = AppendIfMissing(_loginSessions, sess)
				} else if len(sessItems) >= 4 && sessItems[3] == "PRODUCTIVITY" {
					_productivitySessions = AppendIfMissing(_productivitySessions, sess)
				} else {
					_keysToRemove = AppendIfMissing(_keysToRemove, sess)
				}
			}

			totTimeVal, _ := client.Cmd("keys", totTimeEventSearch).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, totTimeVal)

			totCountVal, _ := client.Cmd("keys", totCountEventSearch).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, totCountVal)

			totCountHrVal, _ := client.Cmd("keys", totCountHr).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, totCountHrVal)

			maxTimeVal, _ := client.Cmd("keys", maxTimeEventSearch).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, maxTimeVal)

			thresholdCountVal, _ := client.Cmd("keys", thresholdEventSearch).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, thresholdCountVal)

			thresholdBDCountVal, _ := client.Cmd("keys", thresholdBDEventSearch).List()
			_keysToRemove = AppendListIfMissing(_keysToRemove, thresholdBDCountVal)

		}
		tm := time.Now()
		for _, remove := range _keysToRemove {
			fmt.Println("remove_: ", remove)
			client.Cmd("del", remove)
		}
		for _, session := range _loginSessions {
			fmt.Println("readdSession: ", session)
			client.Cmd("hset", session, "time", tm.Format(layout))
			sessItemsL := strings.Split(session, ":")
			if len(sessItemsL) >= 7 {
				LtotTimeEventName := fmt.Sprintf("TOTALTIME:%s:%s:%s:%s:%s", sessItemsL[1], sessItemsL[2], sessItemsL[3], sessItemsL[5], sessItemsL[6])
				LtotCountEventName := fmt.Sprintf("TOTALCOUNT:%s:%s:%s:%s:%s", sessItemsL[1], sessItemsL[2], sessItemsL[3], sessItemsL[5], sessItemsL[6])
				client.Cmd("set", LtotTimeEventName, 0)
				client.Cmd("set", LtotCountEventName, 0)
			}
		}
		/*for _, prosession := range _productivitySessions {
			fmt.Println("readdSession: ", prosession)
			client.Cmd("hset", prosession, "time", tm.Format(layout))
		}*/
	}
}

func OnSetDailySummary(_date time.Time) {
	totCountEventSearch := fmt.Sprintf("TOTALCOUNT:*")
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

	totalEventKeys, _ := client.Cmd("keys", totCountEventSearch).List()
	for _, key := range totalEventKeys {
		fmt.Println("Key: ", key)
		keyItems := strings.Split(key, ":")
		summery := SummeryDetail{}
		tenant, _ := strconv.Atoi(keyItems[1])
		company, _ := strconv.Atoi(keyItems[2])
		summery.Tenant = tenant
		summery.Company = company
		summery.WindowName = keyItems[3]
		summery.Param1 = keyItems[4]
		summery.Param2 = keyItems[5]

		currentTime := 0
		if summery.WindowName == "LOGIN" {
			sessEventSearch := fmt.Sprintf("SESSION:%d:%d:%s:*:%s:%s", tenant, company, summery.WindowName, summery.Param1, summery.Param2)
			sessEvents, _ := client.Cmd("keys", sessEventSearch).List()
			if len(sessEvents) > 0 {
				tmx, _ := client.Cmd("hget", sessEvents[0], "time").Str()
				tm2, _ := time.Parse(layout, tmx)
				currentTime = int(_date.Sub(tm2.Local()).Seconds())
				fmt.Println("currentTime: ", currentTime)
			}
		}
		totTimeEventName := fmt.Sprintf("TOTALTIME:%d:%d:%s:%s:%s", tenant, company, summery.WindowName, summery.Param1, summery.Param2)
		maxTimeEventName := fmt.Sprintf("MAXTIME:%d:%d:%s:%s:%s", tenant, company, summery.WindowName, summery.Param1, summery.Param2)
		thresholdEventName := fmt.Sprintf("THRESHOLD:%d:%d:%s:%s:%s", tenant, company, summery.WindowName, summery.Param1, summery.Param2)

		fmt.Println("totTimeEventName: ", totTimeEventName)
		fmt.Println("maxTimeEventName: ", maxTimeEventName)
		fmt.Println("thresholdEventName: ", thresholdEventName)

		totCount, _ := client.Cmd("get", key).Int()
		totTime, _ := client.Cmd("get", totTimeEventName).Int()
		maxTime, _ := client.Cmd("get", maxTimeEventName).Int()
		threshold, _ := client.Cmd("get", thresholdEventName).Int()

		fmt.Println("totCount: ", totCount)
		fmt.Println("totTime: ", totTime)
		fmt.Println("maxTime: ", maxTime)
		fmt.Println("threshold: ", threshold)

		summery.TotalCount = totCount
		summery.TotalTime = totTime + currentTime
		summery.MaxTime = maxTime
		summery.ThresholdValue = threshold
		summery.SummaryDate = _date
		go PersistsSummaryData(summery)
	}
}

func OnSetDailyThesholdBreakDown(_date time.Time) {
	thresholdEventSearch := fmt.Sprintf("THRESHOLDBREAKDOWN:*")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnSetDailyThesholdBreakDown", r)
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

	thresholdEventKeys, _ := client.Cmd("keys", thresholdEventSearch).List()
	for _, key := range thresholdEventKeys {
		fmt.Println("Key: ", key)
		keyItems := strings.Split(key, ":")

		if len(keyItems) >= 9 {
			summery := ThresholdBreakDownDetail{}
			tenant, _ := strconv.Atoi(keyItems[1])
			company, _ := strconv.Atoi(keyItems[2])
			hour, _ := strconv.Atoi(keyItems[6])
			summery.Tenant = tenant
			summery.Company = company
			summery.WindowName = keyItems[3]
			summery.Param1 = keyItems[4]
			summery.Param2 = keyItems[5]
			summery.BreakDown = fmt.Sprintf("%s-%s", keyItems[7], keyItems[8])
			summery.Hour = hour

			thCount, _ := client.Cmd("get", key).Int()
			summery.ThresholdCount = thCount
			summery.SummaryDate = _date

			go PersistsThresholdBreakDown(summery)
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

	maxtimeSearch := fmt.Sprintf("MAXTIME:%d:%d:%s:%s:%s", _tenant, _company, _window, _parameter1, _parameter2)
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
			fmt.Println("Recovered in OnGetCurrentMaxTime", r)
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
			timeDiff := int(tm.Local().Sub(tm2.Local()).Seconds())
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
			fmt.Println("Recovered in OnGetCurrentCount", r)
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
		if temptotal < 0 {
			resultChannel <- 0
		} else {
			resultChannel <- temptotal
		}

	} else {
		resultChannel <- 0
	}
}

func OnGetAverageTime(_tenant, _company int, _window, _parameter1, _parameter2 string, resultChannel chan float32) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetAverageTime", r)
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

	tm := time.Now().Local()

	sessEventSearch := fmt.Sprintf("SESSION:%d:%d:%s:*:%s:%s", _tenant, _company, _window, _parameter1, _parameter2)
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

	sessTimeKeyList, _ := client.Cmd("keys", sessEventSearch).List()
	fmt.Println("totalSessTimeKey: ", len(sessTimeKeyList))
	fmt.Println(time.Now().Local())
	if len(sessTimeKeyList) > 0 {
		sessTemptotal := 0
		for _, key := range sessTimeKeyList {
			tmx, _ := client.Cmd("hget", key, "time").Str()
			tm2, _ := time.Parse(layout, tmx)
			timeDiff := int(tm.Local().Sub(tm2.Local()).Seconds())

			if timeDiff > 0 {
				sessTemptotal = sessTemptotal + timeDiff
			}
		}
		totalTime = totalTime + sessTemptotal
	}
	fmt.Println(time.Now().Local())

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

	var avg float32
	if totalCount == 0 {
		avg = 0
	} else {
		avg = float32(totalTime) / float32(totalCount)
	}
	fmt.Println("avg: ", avg)
	resultChannel <- avg
}

func OnGetTotalCount(_tenant, _company int, _window, _parameter1, _parameter2 string, resultChannel chan int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetTotalCount", r)
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

	totalSearch := fmt.Sprintf("TOTALCOUNT:%d:%d:%s:%s:%s", _tenant, _company, _window, _parameter1, _parameter2)
	keyList, _ := client.Cmd("keys", totalSearch).List()
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

func OnGetQueueDetails(_tenant, _company int, resultChannel chan []QueueDetails) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetQueueDetails", r)
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

	queueSearch := fmt.Sprintf("TOTALCOUNT:%d:%d:%s:*", _tenant, _company, "QUEUE")
	keyList, _ := client.Cmd("keys", queueSearch).List()
	if len(keyList) > 0 {
		queueIdList := make([]string, 0)
		for _, key := range keyList {
			keyItems := strings.Split(key, ":")
			if len(keyItems) >= 5 {
				queueIdList = AppendIfMissing(queueIdList, keyItems[4])
			}
		}

		queueDetailList := make([]QueueDetails, 0)

		for _, queueId := range queueIdList {
			queueD := QueueDetails{}

			totalQueued := make(chan int)
			totalAnswer := make(chan int)
			totalDropped := make(chan int)
			maxWaitTime := make(chan int)
			currentMaxWaitTime := make(chan int)
			currentWaiting := make(chan int)
			avgWaitingTime := make(chan float32)

			go OnGetTotalCount(_tenant, _company, "QUEUE", queueId, "*", totalQueued)
			go OnGetTotalCount(_tenant, _company, "QUEUEANSWERED", queueId, "*", totalAnswer)
			go OnGetTotalCount(_tenant, _company, "QUEUEDROPPED", queueId, "*", totalDropped)
			go OnGetMaxTime(_tenant, _company, "QUEUE", queueId, "*", maxWaitTime)
			go OnGetCurrentMaxTime(_tenant, _company, "QUEUE", queueId, "*", currentMaxWaitTime)
			go OnGetCurrentCount(_tenant, _company, "QUEUE", queueId, "*", currentWaiting)
			go OnGetAverageTime(_tenant, _company, "QUEUE", queueId, "*", avgWaitingTime)

			queueD.QueueId = queueId
			queueD.QueueName = GetQueueName(queueId)
			queueD.QueueInfo.TotalQueued = <-totalQueued
			queueD.QueueInfo.TotalAnswered = <-totalAnswer
			queueD.QueueInfo.QueueDropped = <-totalDropped
			queueD.QueueInfo.MaxWaitTime = <-maxWaitTime
			queueD.QueueInfo.CurrentMaxWaitTime = <-currentMaxWaitTime
			queueD.QueueInfo.CurrentWaiting = <-currentWaiting
			queueD.QueueInfo.AverageWaitTime = <-avgWaitingTime

			close(totalQueued)
			close(totalAnswer)
			close(totalDropped)
			close(maxWaitTime)
			close(currentMaxWaitTime)
			close(currentWaiting)
			close(avgWaitingTime)

			queueDetailList = append(queueDetailList, queueD)
		}

		resultChannel <- queueDetailList
	} else {
		resultChannel <- make([]QueueDetails, 0)
	}
}

func OnGetSingleQueueDetails(_tenant, _company int, queueId string, resultChannel chan QueueDetails) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetQueueDetails", r)
		}
	}()

	queueD := QueueDetails{}

	totalQueued := make(chan int)
	totalAnswer := make(chan int)
	totalDropped := make(chan int)
	maxWaitTime := make(chan int)
	currentMaxWaitTime := make(chan int)
	currentWaiting := make(chan int)
	avgWaitingTime := make(chan float32)

	go OnGetTotalCount(_tenant, _company, "QUEUE", queueId, "*", totalQueued)
	go OnGetTotalCount(_tenant, _company, "QUEUEANSWERED", queueId, "*", totalAnswer)
	go OnGetTotalCount(_tenant, _company, "QUEUEDROPPED", queueId, "*", totalDropped)
	go OnGetMaxTime(_tenant, _company, "QUEUE", queueId, "*", maxWaitTime)
	go OnGetCurrentMaxTime(_tenant, _company, "QUEUE", queueId, "*", currentMaxWaitTime)
	go OnGetCurrentCount(_tenant, _company, "QUEUE", queueId, "*", currentWaiting)
	go OnGetAverageTime(_tenant, _company, "QUEUE", queueId, "*", avgWaitingTime)

	queueD.QueueId = queueId
	queueD.QueueName = GetQueueName(queueId)
	queueD.QueueInfo.TotalQueued = <-totalQueued
	queueD.QueueInfo.TotalAnswered = <-totalAnswer
	queueD.QueueInfo.QueueDropped = <-totalDropped
	queueD.QueueInfo.MaxWaitTime = <-maxWaitTime
	queueD.QueueInfo.CurrentMaxWaitTime = <-currentMaxWaitTime
	queueD.QueueInfo.CurrentWaiting = <-currentWaiting
	queueD.QueueInfo.AverageWaitTime = <-avgWaitingTime

	close(totalQueued)
	close(totalAnswer)
	close(totalDropped)
	close(maxWaitTime)
	close(currentMaxWaitTime)
	close(currentWaiting)
	close(avgWaitingTime)

	resultChannel <- queueD
}

func GetQueueName(queueId string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetQueueName", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisPubSubIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()
	//authServer
	authE := client.Cmd("auth", redisPubSubPassword)
	errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", ardsRedisDb)
	errHndlr(r.Err)

	qId := strings.Replace(queueId, "-", ":", -1)
	queueName, _ := client.Cmd("hget", "QueueNameHash", qId).Str()
	fmt.Println("queueName: ", queueName)
	if queueName == "" {
		return queueId
	} else {
		return queueName
	}
}
