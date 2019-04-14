package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type responseInfo struct {
	status   int
	bytes    int64
	duration time.Duration
}

type summaryInfo struct {
	requested int64
	responded int64
}

type serverInfo struct {
	serverHostName string
	serverPort     string
	documentPath   string
}

var (
	requests    *int64
	concurrency *int64
	link        string
	timeOut     *int64
	timeLimit   *int64
)

func main() {

	fmt.Println("Hello from my app")

	requests = flag.Int64("n", 1, "Number of requests to perform")
	concurrency = flag.Int64("c", 1, "Number of multiple requests to make at a time")
	timeOut = flag.Int64("tout", 10, "Seconds to max. wait for each response")
	timeLimit = flag.Int64("tlimit", 20, "Maximum number of seconds to spend for benchmarking")

	flag.Parse()
	link = flag.Arg(0)

	//Check whether input values are valid or not
	flagValidation()

	c := make(chan responseInfo)

	summary := summaryInfo{}

	//Extract the server's info
	getServerInfo()

	//Start the benchmark
	startBenchmarking := time.Now()

	for i := int64(0); i < *concurrency; i++ {
		summary.requested++
		go checkLink(link, c)
	}

	for response := range c {
		if summary.requested < *requests {
			summary.requested++
			go checkLink(link, c)
		}

		summary.responded++
		fmt.Println(response)
		//responseChannel <- response
		if summary.responded == summary.requested {
			break
		}
	}

	stopBenchmarking := time.Now()

	if stopBenchmarking.Sub(startBenchmarking) > time.Duration(int(*timeLimit))*time.Second {
		fmt.Println("Benchmark time out!")
		os.Exit(-1)
	}

}

func checkLink(link string, c chan responseInfo) {
	start := time.Now()
	res, err := http.Get(link)
	end := time.Now()

	//Check timeout waiting for response
	if end.Sub(start) > time.Duration(int(*timeOut))*time.Second {
		fmt.Println("Request time out!")
		os.Exit(-1)
	}
	if err != nil {
		panic(err)
	}
	read, _ := io.Copy(ioutil.Discard, res.Body)

	c <- responseInfo{
		status:   res.StatusCode,
		bytes:    read,
		duration: time.Now().Sub(start),
	}
}

func flagValidation() {
	if flag.NArg() == 0 || link == "" {
		fmt.Println("You must enter a web address.")
		os.Exit(-1)
	}
	if *requests <= 0 {
		fmt.Println("Number of requests to perform must be a positive number. Default is 1.")
		os.Exit(-1)
	}
	if *concurrency <= 0 {
		fmt.Println("Number of multiple requests to make at a time must be a positive number. Default is 1.")
		os.Exit(-1)
	}
	if *timeOut <= 0 || *timeLimit <= 0 {
		fmt.Println("Time out and/or time limit must be a positive number.")
		os.Exit(-1)
	}
	if *requests < *concurrency {
		fmt.Println("Number of requests to perform must be greater than number of multiple requests.")
		os.Exit(-1)
	}
}

func getServerInfo() {
	server := serverInfo{}
	u, err := url.Parse(link)
	if err != nil {
		panic(err)
	}
	server.serverHostName = u.Hostname()
	if len(u.Port()) != 0 {
		server.serverPort = u.Port()
	}
	server.documentPath = u.Path
	fmt.Println("Server info ", server)

}
