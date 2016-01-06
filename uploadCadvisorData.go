package main

import "os"

var CadvisorPort = "18080"

func main() {

	LogRun("sys start")

	pushData()
	iAmAlive()
}

func iAmAlive() {
	f, _ := os.OpenFile("test.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	defer f.Close()

	f.Write([]byte("b"))
}
