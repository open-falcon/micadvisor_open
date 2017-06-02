package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	cpuNum   int64
	countNum int
)

func pushData() {
	cadvisorData, err := getCadvisorData()
	if err != nil {
		LogErr(err, "getcadvisorData err")
		return
	}

	t := time.Now().Unix()
	timestamp := fmt.Sprintf("%d", t)

	cadvDataForOneContainer := strings.Split(cadvisorData, `"aliases":[`)
	for k := 1; k < len(cadvDataForOneContainer); k++ { //Traversal containers and ignore head

		memLimit := getMemLimit(cadvDataForOneContainer[k]) //cadvisor provide the memlimit

		containerId := getContainerId(cadvDataForOneContainer[k]) //cadvisor provide the containerId

		DockerData, _ := getDockerData(containerId) //get container inspect

		endpoint := getEndPoint(DockerData) //there is the hosts file path in the inpect of container

		getCpuNum(DockerData) //we need to give the container CPU ENV

		tag := getTag(DockerData) //recode some other message for a container

		ausge, busge := getUsageData(cadvDataForOneContainer[k]) //get 2 usage because some metric recoding Incremental metric

		cpuuage1 := getBetween(ausge, `"cpu":`, `,"diskio":`)
		cpuuage2 := getBetween(busge, `"cpu":`, `,"diskio":`)
		if err := pushCPU(cpuuage1, cpuuage2, timestamp, tag, containerId, endpoint); err != nil { //get cadvisor data about CPU
			LogErr(err, "pushCPU err in pushData")
		}

		diskiouage := getBetween(ausge, `"diskio":`, `,"memory":`)
		if err := pushDiskIo(diskiouage, timestamp, tag, containerId, endpoint); err != nil { //get cadvisor data about DISKIO
			LogErr(err, "pushDiskIo err in pushData")
		}

		memoryuage := getBetween(ausge, `"memory":`, `,"network":`)
		if err := pushMem(memLimit, memoryuage, timestamp, tag, containerId, endpoint); err != nil { //get cadvisor data about Memery
			LogErr(err, "pushMem err in pushData")
		}

		networkuage1 := getBetween(ausge, `"network":`, `,"task_stats":`)
		networkuage2 := getBetween(busge, `"network":`, `,"task_stats":`)
		if err := pushNet(networkuage1, networkuage2, timestamp, tag, containerId, endpoint); err != nil { //get cadvisor data about net
			LogErr(err, "pushNet err in pushData")
		}
	}
}

func pushCount(metric, usageA, usageB, start, end string, countNum int, timestamp, tags, containerId, endpoint string, weight float64) error {

	temp1, _ := strconv.ParseInt(getBetween(usageA, start, end), 10, 64)
	temp2, _ := strconv.ParseInt(getBetween(usageB, start, end), 10, 64)
	usage := float64(temp2-temp1) / float64(countNum) / weight
	value := fmt.Sprintf("%f", usage)
	if err := pushIt(value, timestamp, metric, tags, containerId, "GAUGE", endpoint); err != nil {
		LogErr(err, "pushIt err in "+metric)
		return err
	}
	return nil
}

func pushNet(networkuage1, networkuage2, timestamp, tags, containerId, endpoint string) error {
	LogRun("pushNet")

	if err := pushCount("net.if.in.bytes", networkuage1, networkuage2, `"rx_bytes":`, `,"rx_packets":`, countNum, timestamp, tags, containerId, endpoint, 1); err != nil {
		return err
	}
	if err := pushCount("net.if.in.packets", networkuage1, networkuage2, `"rx_packets":`, `,"rx_errors":`, countNum, timestamp, tags, containerId, endpoint, 1); err != nil {
		return err
	}
	if err := pushCount("net.if.in.errors", networkuage1, networkuage2, `"rx_errors":`, `,"rx_dropped":`, countNum, timestamp, tags, containerId, endpoint, 1); err != nil {
		return err
	}
	if err := pushCount("net.if.in.dropped", networkuage1, networkuage2, `"rx_dropped":`, `,"tx_bytes":`, countNum, timestamp, tags, containerId, endpoint, 1); err != nil {
		return err
	}
	if err := pushCount("net.if.out.bytes", networkuage1, networkuage2, `"tx_bytes":`, `,"tx_packets":`, countNum, timestamp, tags, containerId, endpoint, 1); err != nil {
		return err
	}
	if err := pushCount("net.if.out.packets", networkuage1, networkuage2, `"tx_packets":`, `,"tx_errors":`, countNum, timestamp, tags, containerId, endpoint, 1); err != nil {
		return err
	}
	if err := pushCount("net.if.out.errors", networkuage1, networkuage2, `"tx_errors":`, `,"tx_dropped":`, countNum, timestamp, tags, containerId, endpoint, 1); err != nil {
		return err
	}
	if err := pushCount("net.if.out.dropped", networkuage1, networkuage2, `"tx_dropped":`, `,"tx_bytes":`, countNum, timestamp, tags, containerId, endpoint, 1); err != nil {
		return err
	}

	return nil
}

func pushMem(memLimit, memoryusage, timestamp, tags, containerId, endpoint string) error {
	LogRun("pushMem")
	memUsageNum := getBetween(memoryusage, `"usage":`, `,"working_set"`)
	fenzi, _ := strconv.ParseInt(memUsageNum, 10, 64)
	fenmu, err := strconv.ParseInt(memLimit, 10, 64)
	if err == nil {
		memUsage := float64(fenzi) / float64(fenmu)
		if err := pushIt(fmt.Sprint(memUsage), timestamp, "mem.memused.percent", tags, containerId, "GAUGE", endpoint); err != nil {
			LogErr(err, "pushIt err in pushMem")
		}
	}
	if err := pushIt(memUsageNum, timestamp, "mem.memused", tags, containerId, "GAUGE", endpoint); err != nil {
		LogErr(err, "pushIt err in pushMem")
	}

	if err := pushIt(fmt.Sprint(fenmu), timestamp, "mem.memtotal", tags, containerId, "GAUGE", endpoint); err != nil {
		LogErr(err, "pushIt err in pushMem")
	}

	memHotUsageNum := getBetween(memoryusage, `"working_set":`, `,"container_data"`)
	fenzi, _ = strconv.ParseInt(memHotUsageNum, 10, 64)
	memHotUsage := float64(fenzi) / float64(fenmu)
	if err := pushIt(fmt.Sprint(memHotUsage), timestamp, "mem.memused.hot", tags, containerId, "GAUGE", endpoint); err != nil {
		LogErr(err, "pushIt err in pushMem")
	}

	return nil
}

func pushDiskIo(diskiouage, timestamp, tags, containerId, endpoint string) error {
	LogRun("pushDiskIo")
	temp := getBetween(diskiouage, `"io_service_bytes":\[`, `,"io_serviced":`)
	readUsage := getBetween(temp, `,"Read":`, `,"Sync"`)

	if err := pushIt(readUsage, timestamp, "disk.io.read_bytes", tags, containerId, "COUNTER", endpoint); err != nil {
		LogErr(err, "pushIt err in pushDiskIo")
	}

	writeUsage := getBetween(temp, `,"Write":`, `}`)

	if err := pushIt(writeUsage, timestamp, "disk.io.write_bytes", tags, containerId, "COUNTER", endpoint); err != nil {
		LogErr(err, "pushIt err in pushDiskIo")
	}

	return nil
}

func pushCPU(cpuuage1, cpuuage2, timestamp, tags, containerId, endpoint string) error {
	LogRun("pushCPU")
	if err := pushCount("cpu.busy", cpuuage1, cpuuage2, `{"total":`, `,"per_cpu_usage":`, countNum, timestamp, tags, containerId, endpoint, 10000000*float64(cpuNum)); err != nil {
		return err
	}

	if err := pushCount("cpu.user", cpuuage1, cpuuage2, `"user":`, `,"sy`, countNum, timestamp, tags, containerId, endpoint, 10000000*float64(cpuNum)); err != nil {
		return err
	}

	if err := pushCount("cpu.system", cpuuage1, cpuuage2, `"system":`, `},`, countNum, timestamp, tags, containerId, endpoint, 10000000*float64(cpuNum)); err != nil {
		return err
	}

	percpu1 := strings.Split(getBetween(cpuuage1, `,"per_cpu_usage":\[`, `\],"user":`), `,`)
	percpu2 := strings.Split(getBetween(cpuuage2, `,"per_cpu_usage":\[`, `\],"user":`), `,`)

	metric := fmt.Sprintf("cpu.core.busy")
	for i, _ := range percpu1 {
		temp1, _ := strconv.ParseInt(percpu1[i], 10, 64)
		temp2, _ := strconv.ParseInt(percpu2[i], 10, 64)
		temp3 := temp2 - temp1
		perCpuUsage := fmt.Sprintf("%f", float64(temp3)/10000000)
		if err := pushIt(perCpuUsage, timestamp, metric, tags+",core="+fmt.Sprint(i), containerId, "GAUGE", endpoint); err != nil {
			LogErr(err, "pushIt err in pushCPU")
			return err
		}
	}
	return nil
}
