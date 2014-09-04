package controllers

import (
	"net/http"
	"encoding/json"

	. "easyTask/serv/models"
	
	"github.com/golang/glog"
)

type monitorData struct {
	Name    string
	Ibuff   int64
	Obuff   int64
	Ncout   int64 // 1
	Icout   int64 // 2
	Scout   int64 // 3
	Ecout   int64 // -1
	Nrec    []string
	Irec    []string
	Erec    []string
	Rappers []string
}

type MonitorHandler struct{}

func (this *MonitorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		this.doGet(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func (this *MonitorHandler) doGet(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	act, ttype, tid := r.FormValue("act"), r.FormValue("ttype"), r.FormValue("tid")

	var res []byte
	switch act {
	case "":
		md, err := this.monitorData()
		if err != nil {
			glog.Errorln(err)
		} else {
			res, _ = json.Marshal(md)
		}
		
		break
	case "info":
		if ttype == "" || tid == "" {
			res = []byte("{}")
			break
		}
		
		tmap, err := InfoData(ttype, tid)
		if err != nil {
			glog.Errorln(err)
			break
		}
		if len(tmap) < 1 {
			res = []byte("{}")
		} else {
			res, _ = json.Marshal(tmap[0])
		}
	}

	w.Write(res)
	
	return
}

func (this *MonitorHandler) monitorData() ([]map[string]interface{}, error) {
	var i int = 0
	mdata := make([]map[string]interface{}, len(TaskTypes))

	for k, v := range TaskTypes {
		si, so := v.BuffSize()
		nc, ic, sc, ec, err := GetCount(k)
		if err != nil {
			return nil, err
		}
		nr, ir, er, err := GetRecord(k)
		if err != nil {
			return nil, err
		}

		one := map[string]interface{}{"Name": k, "Ibuff": si, "Obuff": so, "Ncout": nc, "Icout": ic, "Scout": sc, "Ecout": ec, "Nrec": nr, "Irec": ir, "Erec": er, "Rappers": v.RapperNames()}
		
		mdata[i] = one
		i++
	}

	return mdata, nil
}

