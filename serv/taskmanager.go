package main

import (
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
	"time"
	"runtime"
	"flag"
	"encoding/json"
)

var confJson map[string]interface{}

func init() {
	runtime.GOMAXPROCS(8)

	r, err := os.Open("./app.conf")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&confJson)
	if err != nil {
		panic(err)
	}

	dbInit()
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
	if len(name) > 8 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("name's len < 8"))
	}

	err := createDB(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte("success"))
	
	return
}

func main() {
	flag.Parse()
	fmt.Println("task manager service")
	
	http.HandleFunc("/", HandleRoot)

	http.Handle("/putask", &putTaskHandler{})
	http.Handle("/getask", &getTaskHandler{})
	http.Handle("/uptask", &upTaskHandler{})
	http.HandleFunc("/newtype", handleNewTaskType)
	
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
