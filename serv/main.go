package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"time"

	. "easyTask/serv/controllers"
	. "easyTask/serv/models"
	. "easyTask/serv/common"
)

func init() {
	runtime.GOMAXPROCS(8)	
}

func main() {
	flag.Parse()

	http.Handle("/putask", &PutTaskHandler{})
	http.Handle("/getask", &GetTaskHandler{})
	http.Handle("/uptask", &UpTaskHandler{})
	http.Handle("/sayhi", &SayhiHandler{})
	http.HandleFunc("/newtype", HandleNewTaskType)
	http.HandleFunc("/beat", HandleBeat)

	http.Handle("/monitor", &MonitorHandler{})
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./static/"))))

	s := &http.Server{
		Addr:           ConfJson["listenaddr"].(string),
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	go RapperCleaner()

	fmt.Println("easytask GO...")
	if err := s.ListenAndServe(); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
