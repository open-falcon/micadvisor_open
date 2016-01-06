package main

import (
	"log"
	"os"
	"time"
)

var errLog = "/home/work/uploadCadviosrData/log/err_log_uploadCadvisorData.txt"
var runLog = "/home/work/uploadCadviosrData/log/run_log_uploadCadvisorData.txt"

func LogErr(errInfo interface{}, str string) {
	lf, err := os.OpenFile(errLog, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		os.Exit(1)
	}
	defer lf.Close()

	l := log.New(lf, "", os.O_APPEND)

	l.Println(time.Now().Format("2006-01-02 15:04:05"), errInfo, str)
}

func LogRun(str string) {
	lf, err := os.OpenFile(runLog, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		os.Exit(1)
	}
	defer lf.Close()

	l := log.New(lf, "", os.O_APPEND)

	l.Printf("%s", time.Now().Format("2006-01-02 15:04:05"), str)
}
