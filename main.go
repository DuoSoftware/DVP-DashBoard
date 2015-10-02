// DashBoard project main.go
package main

import (
	"code.google.com/p/gorest"
	"net/http"
)

func main() {

	//fmt.Println("Hello World!")

	InitiateRedis()
	InitiateStatDClient()
	gorest.RegisterService(new(DashBoardEvent))
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(":8787", nil)

	////fmt.Scanln()
	//client, error := goredis.Dial(&goredis.DialConfig{Address: "127.0.0.1:6379"})
	//if error == nil {
	//client.Set("key", "value", 0, 0, false, false)

	//	go OnEvent("3", "1", "CALLSERVER", "CALL", "IVR", "1111111", client)

	//} else {
	//	fmt.Println(error.Error())
	//}
	//fmt.Println("Hello World!")
}
