package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
	"time"
)

func main() {
	// input from bash, at least an empty string
	input := os.Args[1]
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	f1Name := os.Getenv("F1NAME")
	f2Name := os.Getenv("F2NAME")

	// make a resty client
	client := resty.New()

	start1 := time.Now()
	_, err1 := client.R().
		EnableTrace().
		SetBody(input).
		Post(fmt.Sprintf("http://%s:%s/%s", host, port, f1Name))
	elap1 := time.Since(start1)

	if err1 != nil {
		fmt.Printf("error calling %s: %v\n", f1Name, err1)
		return
	}

	start2 := time.Now()
	_, err2 := client.R().
		EnableTrace().
		SetBody(input).
		Post(fmt.Sprintf("http://%s:%s/%s", host, port, f2Name))
	elap2 := time.Since(start2)

	if err2 != nil {
		fmt.Printf("error calling %s: %v\n", f2Name, err2)
		return
	}

	// stdout
	fmt.Printf("%s took %v ms and %s took %v ms\n", f1Name, elap1, f2Name, elap2)
}
