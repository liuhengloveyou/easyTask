package controllers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	. "github.com/liuhengloveyou/easyTask/serv/models"

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
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		res, _ = json.Marshal(md)
	case "info":
		if ttype == "" || tid == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tmap, err := InfoData(ttype, tid)
		if err != nil {
			glog.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if len(tmap) < 1 {
			res = []byte("{}")
		} else {
			tmp, _ := base64.StdEncoding.DecodeString(tmap[0]["info"])
			tmap[0]["info"] = string(tmp)
			if "" != tmap[0]["remark"] {
				tmp, _ = base64.StdEncoding.DecodeString(strings.Replace(tmap[0]["remark"], " ", "", -1))
				tmap[0]["remark"] = string(tmp)
			}

			res, _ = json.Marshal(tmap[0])
		}
	case "del":
		if ttype == "" || tid == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err := UpTaskStatByTid(ttype, tid, -100); err != nil {
			glog.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		res = []byte("OK")
	case "redo":
		if ttype == "" || tid == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err := UpTaskStatByTid(ttype, tid, 1); err != nil {
			glog.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		res = []byte("OK")
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
