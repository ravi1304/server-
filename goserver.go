package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var reqType, x, t int
var noChanges int = 0
var data, allData, req, tym, ipAddress, port, url, name, fname, rNo, cpuPercent, pastCpu string
var start time.Time
var nameoffile, processType string
var reqString string

//var rBody []byte
//var t float64
var allargs []string
var logFile, logFile1 *os.File
var querry int = 0

//var start, elapsed time.Time

func hello(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		start = time.Now()
		for k, v := range r.URL.Query() {
			if k == "reqType" {
				req = strings.Join(v, "")
				reqType, _ = strconv.Atoi(req)
			}
			if k == "data" {
				data = strings.Join(v, "")
			}
			if k == "time" {
				tym = strings.Join(v, "")
				t, _ = strconv.Atoi(tym)
				//t, _ = strconv.ParseFloat(tym, 64)
			}
			if k == "reqNo" {
				rNo = strings.Join(v, "")
			}
			if k == "type" {
				processType = strings.Join(v, "")
			}
			if k == "cpu" {
				cpuPercent = strings.Join(v, "")
				//fmt.Fprintf(logFile, "cpu  %s=", cpuPercent)

			}
			if k != "cpu" {
				querry = 1
			} else {
				querry = 0
			}
		}
		if querry == 1 {
			fmt.Fprintf(logFile, "\n--------------- Request %s started at %s [CPU at start: %s]------------------- \n", rNo, start, cpuPercent)
		}
		w.Header().Add("CPU", cpuPercent)
		//fmt.Fprintf(w, "%s \n", name)
		//w.Header().Set("Content-Type", cpuPercent)
		fmt.Fprintf(w, "%s \ncpu= %s\n", name, cpuPercent)

	case "POST":
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte("Received a POST request\n"))
		//fmt.Fprintf(w, "%s\n", allData)
		data = string(reqBody)
		x = len(data)
		//fmt.Printf("%s\n", reqBody)
		allData = "\nRequest: " + rNo + "\n" + "Request Type: " + strconv.Itoa(reqType) + "\n" + "data: " + data + "  [" + strconv.Itoa(x) + " bytes]" + "\n"
		pastCpu = cpuPercent
		//if r.Method == "POST" && querry == 1 {
		//check if a file name is sent
		//fExist := fileExists(data)
		fmt.Fprintf(logFile, "%s %s\n", r.RemoteAddr, r.URL)
		go startProcess(t, rNo, start, allData, pastCpu)
		//querry = 0
		//}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}

func startProcess(t int, rNo string, start time.Time, allData string, pastCpu string) {
	startTime := start
	//fmt.Fprintf(logFile, "%s\n", allData)
	// if processType == "medium" {
	// 	fmt.Println(fib(33))
	// 	f1, err := os.Create("readfile" + rNo + ".txt")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	if err := f1.Truncate(1e7); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	//zip the file
	// 	nameoffile = "readfile" + rNo + ".txt"
	// 	f, _ := os.Open(nameoffile)
	// 	read := bufio.NewReader(f)
	// 	fdata, _ := ioutil.ReadAll(read)
	// 	nameoffile = strings.Replace(nameoffile, ".txt", ".gz", -1)
	// 	f, _ = os.Create(nameoffile)
	// 	w := gzip.NewWriter(f)
	// 	w.Write(fdata)
	// 	w.Close()
	// 	//remove filez
	// 	e := os.Remove("readfile" + rNo + ".txt")
	// 	if e != nil {
	// 		log.Fatal(e)
	// 	}
	// 	e1 := os.Remove("readfile" + rNo + ".gz")
	// 	if e1 != nil {
	// 		log.Fatal(e)
	//	}
	// } else if processType == "high" {
	// 	fmt.Println(fib(40))
	// } else {
	// 	fmt.Println(fib(25)) //processType == "medium"
	// }

	elapsed := time.Since(startTime)
	pesentCPU, _ := strconv.ParseFloat(cpuPercent, 64)
	pastCPU, _ := strconv.ParseFloat(pastCpu, 64)
	reqCPU := pesentCPU - pastCPU
	fmt.Fprintf(logFile, "\n Request %s took %s [CPU at present: %s, req took: %f, processType: %s]\n", rNo, elapsed, cpuPercent, reqCPU, processType)
	reqString := createString(name, rNo, elapsed)
	if reqString != "" {
		http.Get(reqString) //send the req query to server apprp
	}

}

func fib(n int) int {
	if n <= 1 {
		return n
	} else {
		return fib(n-1) + fib(n-2)
	}
}

func fibo(n int) int {
	if n <= 1 {
		return n
	} else {
		return fibo(n-1) + fibo(n-2)
	}
}

func main() {
	cpuPercent = ""
	logFile1, _ = os.Create("cpustats.txt")
	defer logFile1.Close()
	var err error
	reqType = 999
	data = ""
	t = 0 //in ms
	rNo = ""
	x = len(data)
	//ipAddress = os.Args[3] // comment it when using container
	port = os.Args[1]
	name = os.Args[2]
	//
	// ipAddress = "192.168.0.104"
	// port = "3000"
	// name = "server 1"
	// name = os.Args[1]
	//
	url = ipAddress + port
	fname = "logFile_" + name + ".txt"
	logFile, err = os.Create(fname)
	if err != nil {
		log.Fatal("Log file create:", err)
		return
	}
	defer logFile.Close()

	http.HandleFunc("/", hello)
	fmt.Println("Starting", string(name), "on", url)
	if err := http.ListenAndServe(port, nil); err != nil { //change url to port if using container
		log.Fatal(err)
	}

}

func createString(name string, rNo string, elapsed time.Duration) string {
	// "http://192.168.0.103:8000/?name=server1&rNo=1&elapsedTime=xxx" --final string
	//reqString = sURL + "/?reqType=" + strconv.Itoa(requestid) + "&data=response" + "&reqNo=" + strconv.Itoa(reqNo) + "&time=" + strconv.Itoa(processTime)
	reqString = "http://192.168.0.7:8000" + "/?name=" + name + "&rNo=" + rNo + "&elapsedTime=" + elapsed.String()
	return reqString
}
