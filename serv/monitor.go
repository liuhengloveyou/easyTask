package main

import (
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"

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

type monitorHandler struct{}

func (this *monitorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		this.doGet(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func (this *monitorHandler) doGet(w http.ResponseWriter, r *http.Request) {
	var res []byte
	
	r.ParseForm()
	act, ttype, tid := r.FormValue("act"), r.FormValue("ttype"), r.FormValue("tid")

	switch act {
	case "":
		res, _ = json.Marshal(this.monitorData())
		break
	case "info":
		if ttype == "" || tid == "" {
			res = []byte("{}")
			break
		}
		
		res = this.infoData(ttype, tid)
	}

	w.Write(res)
	
	return
}

func (this *monitorHandler) infoData(ttype, tid string) []byte {
	sqlStr := fmt.Sprintf("SELECT * FROM `tasks_%s` WHERE `tid`=?", ttype)
	rst, err := doQuery(sqlStr, tid)
	if err != nil {
		glog.Errorln(sqlStr, err)
		return nil
	}

	if len(rst) != 1 {
		glog.Errorln(sqlStr, tid, len(rst))
		return nil
	}

	rstStr, _ := json.Marshal(rst[0])
	return rstStr
}

func (this *monitorHandler) monitorData() []monitorData {
	var i int = 0
	mdata := make([]monitorData, len(TaskTypes))

	for k, v := range TaskTypes {
		si, so := v.BuffSize()
		nc, ic, sc, ec := this.getCount(k)
		nr, ir, er := this.getRecord(k)

		one := monitorData{Name: k, Ibuff: si, Obuff: so, Ncout: nc, Icout: ic, Scout: sc, Ecout: ec, Nrec: nr, Irec: ir, Erec: er}

		for k, _ := range v.rappers {
			one.Rappers = append(one.Rappers, k)
		}
		
		mdata[i] = one
		i++
	}

	return mdata
}

func (this *monitorHandler) getCount(ttype string) (ncout, icout, scout, ecout int64) {
	sqlStr := fmt.Sprintf("SELECT `stat`, count(0) cou FROM `tasks_%s` GROUP BY `stat`", ttype)
	rst, err := doQuery(sqlStr)
	if err != nil {
		glog.Errorln(sqlStr, err)
		return
	}

	for _, nv := range rst {
		ri, err := strconv.ParseInt(nv["cou"], 10, 64)
		if err != nil {
			glog.Errorln(err)
			continue
		}

		switch nv["stat"] {
		case "1":
			ncout = ri
		case "2":
			icout = ri
		case "3":
			scout = ri
		case "-1":
			ecout = ri
		}
	}

	return
}

func (this *monitorHandler) getRecord(ttype string) (nr, ir, er []string) {
	const NUM = 20
	sqlStr := fmt.Sprintf("select tid,stat from(select @gid:=@cgid,@cgid:=t1.stat,if(@gid=@cgid,@rank:=@rank+1,@rank:=1) rank,t1.* FROM (select tid,stat from `tasks_%s` ORDER BY `id` DESC) t1,(select @gid:=1,@cgid:=1,@rank:=1) t2)t3 where t3.rank<=%d", ttype, NUM)
	rst, err := doQuery(sqlStr)
	if err != nil {
		glog.Errorln(sqlStr, err)
		return
	}

	for _, nv := range rst {
		switch nv["stat"] {
		case "1":
			nr = append(nr, nv["tid"])
		case "2":
			ir = append(ir, nv["tid"])
		case "-1":
			er = append(er, nv["tid"])
		}
	}
	
	return
}
