package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var Interval time.Duration //检测时间间隔

func main() {
	tmp := os.Getenv("Interval")
	Interval = 60 * time.Second
	tmp1, err := strconv.ParseInt(tmp, 10, 64)
	fmt.Println(tmp1)
	if err == nil {
		Interval = time.Duration(tmp1) * time.Second
	}

	cmd := exec.Command("/home/work/uploadCadviosrData/cadvisor")
	if err = cmd.Start(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("start cadvisor ok", Interval)

	go func() {
		t := time.NewTicker(Interval)
		for {
			<-t.C
			cmd = exec.Command("/home/work/uploadCadviosrData/uploadCadvisorData")
			if err := cmd.Start(); err != nil {
				fmt.Println(err)
				return
			}
			cmd.Wait()
		}
	}()

	for {
		time.Sleep(time.Second * 120)
		if isAlive() {
			clean()
		} else {
			os.Exit(1)
		}
	}

}
func isAlive() bool {
	f, _ := os.OpenFile("test.txt", os.O_CREATE|os.O_APPEND|os.O_RDONLY, 0660)
	defer f.Close()
	read_buf := make([]byte, 32)
	var pos int64 = 0
	n, _ := f.ReadAt(read_buf, pos)
	if n == 0 {
		return false
	}
	return true
}

func clean() {
	f, _ := os.OpenFile("test.txt", os.O_TRUNC, 0660)
	defer f.Close()
}
