package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	. "easyTask/serv/common"
	. "easyTask/serv/controllers"
	. "easyTask/serv/models"
)

func init() {
	runtime.GOMAXPROCS(8)
}

func sigHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		s := <-c
		Sig = "service is suspend ..."
		fmt.Println("Got signal:", s)
	}()
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

	sigHandler()
	go RapperCleaner()

	fmt.Println("easytask GO...", ConfJson["listenaddr"].(string))
	if err := s.ListenAndServe(); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

/*
CREATE DATABASE IF NOT EXISTS `taskmanager` DEFAULT CHARACTER SET utf8;

CREATE TABLE `tasks_demo` (
	`id` int(11) NOT NULL AUTO_INCREMENT,
	`tid` varchar(33) NOT NULL,
	`rid` varchar(32) NOT NULL,
	`info` varchar(1024) NOT NULL,
	`stat` int(11) NOT NULL DEFAULT '0', -- 1 = 新任务; 2 = 正在处理; 3 = 处理成功; -1 = 处理出错
	`addTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`overTime` timestamp NULL DEFAULT '0000-00-00 00:00:00',
	`rapper` varchar(256) DEFAULT NULL,
	`client` varchar(256) DEFAULT NULL,
	`remark` text,
	PRIMARY KEY (`id`),
	UNIQUE KEY `inx_tid` (`tid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/
