package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/golang/glog"
)

type TaskInfo struct {
	Tid  string // 任务ID
	Rid  string // 记录ID
	Info string // 任务内容
}

var (
	confJson  map[string]interface{}
	TaskTypes map[string]*TaskType = make(map[string]*TaskType)
)

func init() {
	runtime.GOMAXPROCS(8)

	r, err := os.Open("./app.conf")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&confJson); err != nil {
		panic(err)
	}

	if err := dbInit(); err != nil {
		panic(err)
	}

	if err := loadTaskType(); err != nil {
		panic(err)
	}
}

func loadTaskType() error {
	tables, err := showTables()
	if err != nil {
		return err
	}

	for i := 0; i < len(tables); i++ {
		if strings.HasPrefix(tables[i], "tasks_") {
			tname := tables[i][6:]
			TaskTypes[tname] = NewTaskType()
			TaskTypes[tname].name = tname
		}
	}

	return nil
}

func GetRapper(ttype, name string) (*TaskType, *Rapper) {
	taskTypeOne, ok := TaskTypes[ttype]
	if ok == false {
		return nil, nil
	}

	rapperOne, ok := taskTypeOne.rappers[name]
	if ok == false {
		return taskTypeOne, nil
	}

	return taskTypeOne, rapperOne
}

func rapperCleaner() {
	for {
		time.Sleep(1 * time.Second)

		for k, v := range TaskTypes {
			for k1, v1 := range (*v).rappers {
				if v1.Beat(false) < 0 {
					continue // 已经死了就不要打扰人家
				}
				if (time.Now().Unix() - v1.Beat(false)) > int64(confJson["RapperBeatOut"].(float64)) {
					v.resetRapper(v1)
					v1.Kill()
					glog.Warningln("KILL RAPPER: ", k, k1)
				}
			}
		}

	}
}

func handleNewTaskType(w http.ResponseWriter, r *http.Request) {
	const USAGE = "GET /newtype?name=typename"

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	name := r.FormValue("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(USAGE))
		return

	} else if len(name) > 8 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("name's len < 8"))
		return
	}

	err := createDB(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	TaskTypes[name] = NewTaskType()
	TaskTypes[name].name = name

	w.Write([]byte("OK"))

	return
}

func handleBeat(w http.ResponseWriter, r *http.Request) {
	const USAGE = "GET /beat?type=tasktype&name=rappername"

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	ttype, name := r.FormValue("type"), r.FormValue("name")
	if ttype == "" || name == "" {
		glog.Errorln("beat ERR:", ttype, name)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(USAGE))
		return
	}

	taskTypeOne, rapperOne := GetRapper(ttype, name)
	if taskTypeOne == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no such task type."))
		return
	} else if rapperOne == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no such rapper"))
		return
	}

	rapperOne.Beat(true)

	w.Write([]byte("OK"))
	return
}

func main() {
	flag.Parse()

	http.Handle("/putask", &putTaskHandler{})
	http.Handle("/getask", &getTaskHandler{})
	http.Handle("/uptask", &upTaskHandler{})
	http.Handle("/sayhi", &sayhiHandler{})
	http.HandleFunc("/newtype", handleNewTaskType)
	http.HandleFunc("/beat", handleBeat)

	http.Handle("/monitor", &monitorHandler{})
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./static/"))))

	s := &http.Server{
		Addr:           confJson["listenaddr"].(string),
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	go rapperCleaner()

	fmt.Println("easytask GO...")
	if err := s.ListenAndServe(); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
