package main

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func getCpuNum(dockerdata string) {
	cpuNum = 1
	tmp := getBetween(dockerdata, `"CPU=`, `",`)
	if tmp != "" {
		cpuNum, _ = strconv.ParseInt(tmp, 10, 32)
		if cpuNum == 0 {
			cpuNum = 1
		}
	}
}

func getTag(DockerData string) string {
	//FIXMI:some other message for container
	tags := getBetween(DockerData, `"Tags=`, `",`)
	if tags != "" {
		return tags
	}
	return ""
}

func getMemLimit(str string) string {
	return getBetween(str, `"memory":{"limit":`, `,"`)
}

func getBetween(str, start, end string) string {
	res := regexp.MustCompile(start + `(.+?)` + end).FindStringSubmatch(str)
	if len(res) <= 1 {
		LogErr(errors.New("regexp len < 1"), start+" "+end)
		return ""
	}
	return res[1]
}

func getCadvisorData() (string, error) {
	var (
		resp *http.Response
		err  error
		body []byte
	)
	url := "http://localhost:" + CadvisorPort + "/api/v1.2/docker"
	if resp, err = http.Get(url); err != nil {
		LogErr(err, "Get err in getCadvisorData")
		return "", err
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		LogErr(err, "ReadAll err in getCadvisorData")
		return "", err
	}

	return string(body), nil
}

func getUsageData(cadvisorData string) (ausge, busge string) {
	ausge = strings.Split(cadvisorData, `{"timestamp":`)[1]
	if len(strings.Split(cadvisorData, `{"timestamp":`)) < 11 {
		countNum = 1
		busge = strings.Split(cadvisorData, `{"timestamp":`)[2]
	} else {
		busge = strings.Split(cadvisorData, `{"timestamp":`)[11]
		countNum = 10
	}

	return ausge, busge
}

func getContainerId(cadvisorData string) string {

	getContainerId1 := strings.Split(cadvisorData, `],"namespace"`)
	getContainerId2 := strings.Split(getContainerId1[0], `","`)
	getContainerId3 := strings.Split(getContainerId2[1], `"`)
	containerId := getContainerId3[0]

	return containerId
}

func getEndPoint(DockerData string) string {
	//get endpoint from env first
	endPoint := getBetween(DockerData, `"EndPoint=`, `",`)
	if endPoint != "" {
		return endPoint
	}
	//get docker name
	docker_name := getBetween(DockerData, `"Name":"`, `",`)
	if docker_name != "" {
		return docker_name
	}
	filepath := getBetween(DockerData, `"HostsPath":"`, `",`)
	buf := make(map[int]string, 6)
	inputFile, inputError := os.Open(filepath)
	if inputError != nil {
		LogErr(inputError, "getEndPoint open file err"+filepath)
		return ""
	}
	defer inputFile.Close()

	inputReader := bufio.NewReader(inputFile)
	lineCounter := 0
	for i := 0; i < 2; i++ {
		inputString, readerError := inputReader.ReadString('\n')
		if readerError == io.EOF {
			break
		}
		lineCounter++
		buf[lineCounter] = inputString
	}
	hostname := strings.Split(buf[1], "	")[0]
	hostname = strings.Replace(hostname, "\n", " ", -1)
	return hostname
}

func getDockerData(containerId string) (string, error) {
	str, err := RequestUnixSocket("/containers/"+containerId+"/json", "GET")
	if err != nil {
		LogErr(err, "getDockerData err")
	}
	return str, nil
}

func RequestUnixSocket(address, method string) (string, error) {
	DOCKER_UNIX_SOCKET := "unix:///var/run/docker.sock"
	// Example: unix:///var/run/docker.sock:/images/json?since=1374067924
	unix_socket_url := DOCKER_UNIX_SOCKET + ":" + address
	u, err := url.Parse(unix_socket_url)
	if err != nil || u.Scheme != "unix" {
		LogErr(err, "Error to parse unix socket url "+unix_socket_url)
		return "", err
	}

	hostPath := strings.Split(u.Path, ":")
	u.Host = hostPath[0]
	u.Path = hostPath[1]

	conn, err := net.Dial("unix", u.Host)
	if err != nil {
		LogErr(err, "Error to connect to"+u.Host)
		// fmt.Println("Error to connect to", u.Host, err)
		return "", err
	}

	reader := strings.NewReader("")
	query := ""
	if len(u.RawQuery) > 0 {
		query = "?" + u.RawQuery
	}

	request, err := http.NewRequest(method, u.Path+query, reader)
	if err != nil {
		LogErr(err, "Error to create http request")
		// fmt.Println("Error to create http request", err)
		return "", err
	}

	client := httputil.NewClientConn(conn, nil)
	response, err := client.Do(request)
	if err != nil {
		LogErr(err, "Error to achieve http request over unix socket")
		// fmt.Println("Error to achieve http request over unix socket", err)
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		LogErr(err, "Error, get invalid body in answer")
		// fmt.Println("Error, get invalid body in answer")
		return "", err
	}

	defer response.Body.Close()

	return string(body), err
}
