package main

import (
	"flag"
	"fmt"
)

func main() {
	fmt.Println("Hello from my app")

	requests := flag.Int64("n", 1, "Number of requests to perform")
	concurrency := flag.Int64("c", 1, "Number of multiple requests to make at a time")

	fmt.Println(&requests, &concurrency)
	fmt.Println(*requests, *concurrency)
	flag.Parse()
	flag.PrintDefaults()
}
