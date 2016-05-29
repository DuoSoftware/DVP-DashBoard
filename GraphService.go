package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func IncokeGhaphite(_url string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in IncokeGhaphite", r)
		}
	}()
	fmt.Println("graphite_url: ", _url)
	resp, err := http.Get(_url)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	} else {
		defer resp.Body.Close()
		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			return ""
		} else {
			tmx := string(response[:])
			fmt.Println(tmx)
			return tmx
		}
	}
}

func OnGetCalls(_tenant, _company, _duration int) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetCalls", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=summarize(derivative(stats.event.concurrent.%d.%d.*.CALLS),\"1min\")&from=-%dmin&format=json", statsDIp, _tenant, _company, _duration)
	return IncokeGhaphite(url)
}

func OnGetChannels(_tenant, _company, _duration int) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetChannels", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=summarize(derivative(stats.event.concurrent.%d.%d.*.CALLCHANNELS),\"1min\")&from=-%dmin&format=json", statsDIp, _tenant, _company, _duration)
	return IncokeGhaphite(url)
}

func OnGetBridge(_tenant, _company, _duration int) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetBridge", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=summarize(derivative(stats.event.concurrent.%d.%d.*.BRIDGE),\"1min\")&from=-%dmin&format=json", statsDIp, _tenant, _company, _duration)
	return IncokeGhaphite(url)
}

func OnGetQueued(_tenant, _company, _duration int) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetQueued", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=summarize(derivative(stats.event.concurrent.%d.%d.*.QUEUE),\"1min\")&from=-%dmin&format=json", statsDIp, _tenant, _company, _duration)
	return IncokeGhaphite(url)
}
