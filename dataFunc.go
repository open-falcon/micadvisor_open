package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func pushIt(value, timestamp, metric, tags, containerId, counterType, endpoint string) error {

	postThing := `[{"metric": "` + metric + `", "endpoint": "` + endpoint + `", "timestamp": ` + timestamp + `,"step": ` + "60" + `,"value": ` + value + `,"counterType": "` + counterType + `","tags": "` + tags + `"}]`
	LogRun(postThing)
	//push data to falcon-agent
	url := "http://127.0.0.1:1988/v1/push"
	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(postThing))
	if err != nil {
		LogErr(err, "Post err in pushIt")
		return err
	}
	defer resp.Body.Close()
	_, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		LogErr(err1, "ReadAll err in pushIt")
		return err1
	}
	return nil
}
