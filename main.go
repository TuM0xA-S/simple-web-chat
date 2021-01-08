package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

type message struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
	Time   string `json:"time"`
}

var messages = []message{}

var mutex sync.Mutex

func sendHandler(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return
	}
	fmt.Println(ip, "is sending message")
	bytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		return
	}

	text := string(bytes)
	time := time.Now().Format("15:04")

	msg := message{
		Sender: ip,
		Text:   text,
		Time:   time,
	}

	mutex.Lock()
	messages = append(messages, msg)
	mutex.Unlock()
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var cnt int
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	delta := []message{}
	if err := dec.Decode(&cnt); err == nil && cnt >= 0 && cnt < len(messages) {
		delta = messages[cnt:]
	}

	bytes, _ := json.Marshal(delta)
	w.Write(bytes)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/update", updateHandler)

	var host string
	args := os.Args[1:]
	if len(args) == 0 {
		host = ":80"
	} else if len(args) == 1 {
		host = args[0]
	} else {
		log.Fatal("usage: chat-server [host]")
	}
	if err := http.ListenAndServe(host, nil); err != nil {
		log.Fatal(err)
	}
}
