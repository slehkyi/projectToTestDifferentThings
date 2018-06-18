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
const srv1 = "http://127.0.0.1:8080/check"
const srv2 = "http://127.0.0.1:8088/check"

func main() {

	go Change(nil)
	mux := http.NewServeMux()
	mux.HandleFunc("/check", CheckHandler)
	go http.ListenAndServe(":8080", mux)

	for {
		time.Sleep(2 * time.Second)
		go Change(mux)
	}
}

func Change(mux *http.ServeMux) {
	result, err := Check(srv2)
	if err != nil {
		Active(true)
		fmt.Println("I am Active")
		fmt.Println(err)
		mux = http.NewServeMux()
		mux.HandleFunc("/check", CheckHandler)
		http.ListenAndServe(":8080", mux)
	}
	if result.Active == false {
		fmt.Println("I am Active")
		Active(true)
		fmt.Println(AmActive)
	} else {
		fmt.Println("I am not active, but listening")
		Active(false)
		fmt.Println(AmActive)
	}
}

type responseOK struct {
	Active bool
	TimeStamp int64
}

func CheckHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	bodyActive := responseOK{
		Active: AmActive,
		TimeStamp: time.Now().UnixNano(),
	}
	j, err := json.Marshal(bodyActive)
	if err != nil {
		log.Println(err)
	}
	res.Write(j)
}

func Check(url string) (responseOK, error) {
	var result responseOK

	resp, err := http.Get(url)
	if err != nil {
		return responseOK{}, err
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

	return result, nil
}

func Active(val bool) {
	AmActive = val
}
