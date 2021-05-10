package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func PersistSessionInfo(tenant, company int, businessUnit, window, dashboardSession, param1, param2, timeStr string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in PersistSessionInfo", r)
		}
	}()

	dialUrl := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", mongoUser, mongoPassword, mongoIp, mongoPort, mongoDbname)
	session, err := mgo.Dial(dialUrl)

	if err != nil {
		//panic(err)
		log.Println(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB(mongoDbname).C("DashboardSessions")

	index := mgo.Index{
		Key:        []string{"Tenant", "Company", "Window", "Session"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	}
	errIndex := c.EnsureIndex(index)
	if errIndex != nil {
		log.Fatal(errIndex)
	}

	err = c.Insert(&SessionPersistence{tenant, company, businessUnit, window, dashboardSession, param1, param2, timeStr})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Persist Session Successful:: ", dashboardSession)
	}
}

func FindPersistedSession(tenant, company int, window, dashboardSession string) (sessionKey, timeValue, businessUnit, param1, param2 string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in FindPersistedSession", r)
		}
	}()

	dialUrl := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", mongoUser, mongoPassword, mongoIp, mongoPort, mongoDbname)
	session, err := mgo.Dial(dialUrl)

	if err != nil {
		//panic(err)
		log.Println(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB(mongoDbname).C("DashboardSessions")

	result := SessionPersistence{}
	err = c.Find(bson.M{"tenant": tenant, "company": company, "window": window, "session": dashboardSession}).One(&result)
	if err != nil {
		fmt.Println(err)
	} else {

		fmt.Println("SessionPersistenceInfo :", result)
		timeValue = result.Time
		sessionKey = fmt.Sprintf("SESSION:%d:%d:%s:%s:%s:%s:%s", result.Tenant, result.Company, result.BusinessUnit, result.Window, result.Session, result.Param1, result.Param2)
		param1 = result.Param1
		param2 = result.Param2
		businessUnit = result.BusinessUnit
	}
	return
}

func DeletePersistedSession(tenant, company int, window, dashboardSession string) (result int64) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in DeletePersistedSession", r)
		}
	}()
	result = 0
	dialUrl := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", mongoUser, mongoPassword, mongoIp, mongoPort, mongoDbname)
	session, err := mgo.Dial(dialUrl)

	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB(mongoDbname).C("DashboardSessions")

	//result := SessionPersistence{}
	err = c.Remove(bson.M{"tenant": tenant, "company": company, "window": window, "session": dashboardSession})
	if err != nil {
		fmt.Println(err)
	} else {

		fmt.Println("Remove SessionPersistenceInfo :Success")
		result = 1
	}
	return
}
