package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var logFile1 *os.File
var serverName, cpu, url, cpuper string

func cpuStats() {
	out, err := exec.Command("/usr/bin/docker", "stats", "--no-stream", "--format", "table {{.Name}}\t{{.CPUPerc}}").Output()
	if err != nil {
		fmt.Fprintf(logFile1, "error %s \n", err)
	} else {

		output := string(out[:])
		//split the output by words
		f := strings.Fields(output)
		// iterate over array of words
		if len(f) > 3 { //if only i have some output then print inside the log file else there is no docker running
			fmt.Fprintf(logFile1, "Command Successfully Executed \n")
			fmt.Fprintf(logFile1, "%s\n", output)
		} else {
			fmt.Println("command executed but no docker container is runnng") //to show that the command is ececutes properly
		}
		for i := 3; i < len(f); i = i + 2 {
			serverName = f[i]
			cpu = f[i+1]
			split := strings.Split(cpu, "%")
			fmt.Println(split)

			//fmt.Println(serverName, "cpu=", cpu, "split", split[0])
			cpuper = split[0]
			go sendToServers(serverName, cpuper)

		}
	}
}

func sendToServers(serverName string, cpuper string) {
	if serverName == "server1" {
		url = "http://192.168.0.4:5000/?cpu="
	} else if serverName == "server2" {
		url = "http://192.168.0.3:5000/?cpu="
	} else if serverName == "server3" {
		url = "http://192.168.0.108:5000/?cpu="
	} else if serverName == "server4" {
		url = "http://192.168.0.106:5000/?cpu="
	} else if serverName == "server5" {
		url = "http://192.168.0.107:5000/?cpu="
	} else {
		fmt.Println("server name did not match")
	}
	cpuReqString := url + cpuper
	//fmt.Println(cpuReqString)
	out1, err := exec.Command("curl", cpuReqString).Output()
	if err != nil {
		fmt.Printf("error %s \n", err)
	} else {
		output1 := string(out1[:])
		fmt.Println("response from server", output1)
	}

}

func main() {
	logFile1, _ = os.Create("cpustats.txt")
	defer logFile1.Close()
	if runtime.GOOS == "windows" {
		fmt.Println("can't execute on windows machine")

	} else {
		for {
			//<-time.Tick(1 * time.Second)
			<-time.Tick(10 * time.Millisecond)
			cpuStats()
		}
	}
}
