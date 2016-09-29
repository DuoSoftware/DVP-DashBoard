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
	//resp, err := http.Get(_url)
	//resp, err := http.NewRequest("GET", _url, nil)

	tr := &http.Transport{
		DisableCompression: true,
		DisableKeepAlives:  true,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(_url)

	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		result <- ""
	} else {
		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			result <- ""
		} else {
			tmx := string(response[:])
			fmt.Println(tmx)
			result <- tmx
		}
	}
}

func OnGetCalls(_tenant, _company, _duration int, result chan string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetCalls", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=stats.event.common.concurrent.%d.%d.*.CALLS&from=-%dmin&format=json", statsDIp, _tenant, _company, _duration)
	resultChannel := make(chan string)
	go IncokeGhaphite(url, resultChannel)
	var queueInfo = <-resultChannel
	close(resultChannel)
	result <- queueInfo
}

func OnGetChannels(_tenant, _company, _duration int, result chan string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetChannels", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=stats.event.common.concurrent.%d.%d.*.CALLCHANNELS&from=-%dmin&format=json", statsDIp, _tenant, _company, _duration)
	resultChannel := make(chan string)
	go IncokeGhaphite(url, resultChannel)
	var queueInfo = <-resultChannel
	close(resultChannel)
	result <- queueInfo
}

func OnGetBridge(_tenant, _company, _duration int, result chan string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetBridge", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=stats.event.common.concurrent.%d.%d.*.BRIDGE&from=-%dmin&format=json", statsDIp, _tenant, _company, _duration)
	resultChannel := make(chan string)
	go IncokeGhaphite(url, resultChannel)
	var queueInfo = <-resultChannel
	close(resultChannel)
	result <- queueInfo
}

func OnGetQueued(_tenant, _company, _duration int, result chan string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetQueued", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=sumSeries(stats.event.common.concurrent.%d.%d.*.QUEUE)&from=-%dmin&format=json", statsDIp, _tenant, _company, _duration)
	resultChannel := make(chan string)
	go IncokeGhaphite(url, resultChannel)
	var queueInfo = <-resultChannel
	close(resultChannel)
	result <- queueInfo
}

func OnGetConcurrentQueue(_tenant, _company, _duration int, _queue string, result chan string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetQueued", r)
		}
	}()

	//stats.gauges.event.concurrent.1.3.Queue-3-1-CALLSERVER-CALL-attribute_8-L.QUEUE
	url := fmt.Sprintf("http://%s/render?target=stats.gauges.event.common.concurrent.%d.%d.%s.QUEUE&from=-%dmin&format=json", statsDIp, _tenant, _company, _queue, _duration)
	resultChannel := make(chan string)
	go IncokeGhaphite(url, resultChannel)
	var queueInfo = <-resultChannel
	close(resultChannel)
	result <- queueInfo
}

func OnGetTotalNewTicket(_tenant, _company, _duration int, result chan string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetTotalNewTicket", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=summarize(sumSeries(stats.gauges.event.ticket.totalcount.%d.%d.total.NEWTICKET),\"1d\",\"max\",true)&from=-%dd&format=json", statsDIp, _tenant, _company, _duration)
	resultChannel := make(chan string)
	go IncokeGhaphite(url, resultChannel)
	var ticketInfo = <-resultChannel
	close(resultChannel)
	result <- ticketInfo
}

func OnGetTotalClosedTicket(_tenant, _company, _duration int, result chan string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetTotalClosedTicket", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=summarize(sumSeries(stats.gauges.event.ticket.totalcount.%d.%d.total.CLOSEDTICKET),\"1d\",\"max\",true)&from=-%dd&format=json", statsDIp, _tenant, _company, _duration)
	resultChannel := make(chan string)
	go IncokeGhaphite(url, resultChannel)
	var ticketInfo = <-resultChannel
	close(resultChannel)
	result <- ticketInfo
}

func OnGetDiffClosedVsNew(_tenant, _company, _duration int, result chan string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetDiffClosedVsNew", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=summarize(diffSeries(stats.gauges.event.ticket.totalcount.%d.%d.total.NEWTICKET,stats.gauges.event.totalcount.%d.%d.total.CLOSEDTICKET),\"1d\",\"max\",true)&from=-%dd&format=json", statsDIp, _tenant, _company, _tenant, _company, _duration)
	resultChannel := make(chan string)
	go IncokeGhaphite(url, resultChannel)
	var ticketInfo = <-resultChannel
	close(resultChannel)
	result <- ticketInfo
}

func OnGetTotalNewTicketByUser(_tenant, _company, _duration int, _userName string, result chan string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetTotalNewTicket", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=summarize(sumSeries(stats.gauges.event.ticket.totalcount.%d.%d.user_%s.NEWTICKET),\"1d\",\"max\",true)&from=-%dd&format=json", statsDIp, _tenant, _company, _userName, _duration)
	resultChannel := make(chan string)
	go IncokeGhaphite(url, resultChannel)
	var ticketInfo = <-resultChannel
	close(resultChannel)
	result <- ticketInfo
}

func OnGetTotalClosedTicketByUser(_tenant, _company, _duration int, _userName string, result chan string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetTotalClosedTicket", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=summarize(sumSeries(stats.gauges.event.ticket.totalcount.%d.%d.user_%s.CLOSEDTICKET),\"1d\",\"max\",true)&from=-%dd&format=json", statsDIp, _tenant, _company, _userName, _duration)
	resultChannel := make(chan string)
	go IncokeGhaphite(url, resultChannel)
	var ticketInfo = <-resultChannel
	close(resultChannel)
	result <- ticketInfo
}

func OnGetDiffClosedVsNewByUser(_tenant, _company, _duration int, _userName string, result chan string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in OnGetDiffClosedVsNew", r)
		}
	}()
	url := fmt.Sprintf("http://%s/render?target=summarize(diffSeries(stats.gauges.event.ticket.totalcount.%d.%d.user_%s.NEWTICKET,stats.gauges.event.totalcount.%d.%d.user_%s.CLOSEDTICKET),\"1d\",\"max\",true)&from=-%dd&format=json", statsDIp, _tenant, _company, _userName, _tenant, _company, _userName, _duration)
	resultChannel := make(chan string)
	go IncokeGhaphite(url, resultChannel)
	var ticketInfo = <-resultChannel
	close(resultChannel)
	result <- ticketInfo
}
