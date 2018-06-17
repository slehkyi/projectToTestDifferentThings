package main

import (
	"net/http"
	"encoding/json"
	"log"
	"time"
	"fmt"
	"io/ioutil"
)

var AmActive = false

func main() {

	var result responseOK

	resp, err := http.Get("http://127.0.0.1:8088/check")
	if err != nil {
		http.ListenAndServe(":8080", nil)
		Active(true)
		fmt.Println("I am Active")
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

	http.HandleFunc("/check", Check)
}

type responseOK struct {
	status string
	timeStamp int64
}

func Check(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	body := responseOK{
		status: "OK",
		timeStamp: time.Now().UnixNano(),
	}
	if AmActive == true {
		j, err := json.Marshal(body)
		if err != nil {
			log.Println(err)
		}
		res.Write(j)
	}
}

func Active(val bool) {
	AmActive = val
}
