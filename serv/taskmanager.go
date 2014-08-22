package main

import (
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"time"
	"runtime"
	"flag"
	"encoding/json"
)

type TaskInfo struct {
	Tid string // 任务ID
	Rid string // 记录ID
	Info string // 任务内容
}

var (
	confJson map[string]interface{}
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

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	fn := r.URL.Path[1:]
	if fn == "" {
		w.Write([]byte("task manager service\n"))
		return
	}

	contents, err := ioutil.ReadFile("./web/" + fn)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write(contents)

	return
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("%v", time.Now().UnixNano())))
}

/*
GET /newtype?name=typename
*/
func handleNewTaskType(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	name := r.FormValue("name")
	if name == "" || len(name) > 8 {
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

func main() {
	flag.Parse()
	fmt.Println("task manager service")
	
	http.HandleFunc("/", HandleRoot)

	http.Handle("/putask", &putTaskHandler{})
	http.Handle("/getask", &getTaskHandler{})
	http.Handle("/uptask", &upTaskHandler{})
	http.Handle("/sayhi", &handleSayhi{})
	http.HandleFunc("/newtype", handleNewTaskType)
	http.HandleFunc("/beat", handleNewTaskType)
	
	http.HandleFunc("/upload", handleUpload)
	
	s := &http.Server{
		Addr:          confJson["listenaddr"].(string),
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
