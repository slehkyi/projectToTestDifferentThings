package main

import (
	"net/http"
	"encoding/json"
	"log"
	"time"
	"fmt"
	"io/ioutil"
	"net"
)

var AmActive = false
var urlCheck string
const IP1 = "192.168.0.214"
const IP2 = "127.0.0.1"
const port1 = ":8080"
const port2 = ":8088"
const urlCheck1 = "http://"+IP1+port1+"/check"
const urlCheck2 = "http://"+IP2+port2+"/check"

func main() {

	myIP := GetOutboundIP()
	fmt.Println(myIP.String())

	// go Change(nil)
	mux := http.NewServeMux()
	mux.HandleFunc("/check", CheckHandler)
	go selectServer(myIP.String(), mux)

	for {
		time.Sleep(2 * time.Second)
		go Change(mux)
	}
}

func selectServer(myIP string, mux *http.ServeMux) {
	switch myIP {
	case IP1:
		http.ListenAndServe(port1, mux)
		urlCheck = "http://"+IP2+port2+"/check"
	case IP2:
		http.ListenAndServe(port2, mux)
		urlCheck = "http://"+IP1+port1+"/check"
	}
}

func Change(mux *http.ServeMux) {
	result, err := Check(urlCheck2)
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

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func Active(val bool) {
	AmActive = val
}
