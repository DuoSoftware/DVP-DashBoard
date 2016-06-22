package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func IncokeGhaphite(_url string, result chan string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in IncokeGhaphite", r)
		}
	}()
	fmt.Println("graphite_url: ", _url)
	resp, err := http.Get(_url)

	if err != nil {
		fmt.Println(err.Error())
		resp.Body.Close()
		result <- ""
	} else {
		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			resp.Body.Close()
			result <- ""
		} else {
			tmx := string(response[:])
			fmt.Println(tmx)
			resp.Body.Close()
			result <- tmx
		}
	}
}

func OnGetCalls(_tenant, _company, _duration int) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetCalls", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=stats.event.concurrent.%d.%d.*.CALLS&from=-%dmin&format=json", statsDIp, _tenant, _company, _duration)
	resultChannel := make(chan string)
	go IncokeGhaphite(url, resultChannel)
	var queueInfo = <-resultChannel
	close(resultChannel)
	return queueInfo
}

func OnGetChannels(_tenant, _company, _duration int) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetChannels", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=stats.event.concurrent.%d.%d.*.CALLCHANNELS&from=-%dmin&format=json", statsDIp, _tenant, _company, _duration)
	resultChannel := make(chan string)
	go IncokeGhaphite(url, resultChannel)
	var queueInfo = <-resultChannel
	close(resultChannel)
	return queueInfo
}

func OnGetBridge(_tenant, _company, _duration int) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetBridge", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=stats.event.concurrent.%d.%d.*.BRIDGE&from=-%dmin&format=json", statsDIp, _tenant, _company, _duration)
	resultChannel := make(chan string)
	go IncokeGhaphite(url, resultChannel)
	var queueInfo = <-resultChannel
	close(resultChannel)
	return queueInfo
}

func OnGetQueued(_tenant, _company, _duration int) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetQueued", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=stats.event.concurrent.%d.%d.*.QUEUE&from=-%dmin&format=json", statsDIp, _tenant, _company, _duration)
	resultChannel := make(chan string)
	go IncokeGhaphite(url, resultChannel)
	var queueInfo = <-resultChannel
	close(resultChannel)
	return queueInfo
}

func OnGetConcurrentQueue(_tenant, _company, _duration int, _queue string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetQueued", r)
		}
	}()

	//stats.gauges.event.concurrent.1.3.Queue-3-1-CALLSERVER-CALL-attribute_8-L.QUEUE
	url := fmt.Sprintf("http://%s/render?target=stats.gauges.event.concurrent.%d.%d.%s.QUEUE&from=-%dmin&format=json", statsDIp, _tenant, _company, _queue, _duration)
	resultChannel := make(chan string)
	go IncokeGhaphite(url, resultChannel)
	var queueInfo = <-resultChannel
	close(resultChannel)
	return queueInfo
}
