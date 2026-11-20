package main

import (
	"cancer"
	"cancer/log"
)

var logger = log.RegisterScope("/main")

func main() {

	server := cancer.NewAppServer()
	var port int = 9999
	server.Bind(port)
	server.Accept("/")

}
