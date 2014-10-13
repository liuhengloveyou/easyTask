package models

import (
	"fmt"
	"strconv"
)

func CreateDB(name string) (err error) {
	const TABLESQL = "CREATE TABLE `tasks_%s` (`id` int(11) NOT NULL AUTO_INCREMENT,`tid` varchar(33) NOT NULL,`rid` varchar(32) NOT NULL,`info` varchar(1024) NOT NULL,`stat` int(11) NOT NULL DEFAULT '0',`addTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,`overTime` timestamp NULL DEFAULT '0000-00-00 00:00:00',`rapper` varchar(256) DEFAULT NULL,`client` varchar(256) DEFAULT NULL,`remark` text,PRIMARY KEY (`id`),UNIQUE KEY `inx_tid` (`tid`)) ENGINE=InnoDB DEFAULT CHARSET=utf8"

	sqlStr := fmt.Sprintf(TABLESQL, name)
	_, err = db.Insert(sqlStr)

	return
}

func ShowTables() ([]string, error) {
	tables, err := db.Query("show tables")
	if err != nil {
		return nil, err
	}

	rst := make([]string, len(tables))
	for n, tn := range tables {
		rst[n] = tn["Tables_in_taskmanager"]
	}

	return rst, nil
}

func newTask2DB(ttype, tid, rid, info string, stat int64) (int64, error) {
	sqlStr := fmt.Sprintf("INSERT INTO `tasks_%s`(`tid`, `rid`, `info`, `stat`) VALUES(?,?,?,?)", ttype)
	return db.Insert(sqlStr, tid, rid, info, stat)
}

func upTask2DB(ttype, tid, rapper, msg string, stat int64) (int64, error) {
	sqlStr := fmt.Sprintf("UPDATE `tasks_%s` SET `stat`=?, `overTime`=CURRENT_TIMESTAMP, `rapper`=?, `remark`=? WHERE `tid`=?", ttype)
	return db.Update(sqlStr, stat, rapper, msg, tid)
}

func upTaskStat2DB(ttype string, ids, ide int64) (int64, error) {
	sqlStr := fmt.Sprintf("UPDATE `tasks_%s` SET `stat`=2 WHERE `id` >= ? AND `id` <= ? AND `stat`=1", ttype)
	return db.Update(sqlStr, ide, ids)
}

func getTasks(ttype string, num int64) ([]map[string]string, error) {
	sqlStr := fmt.Sprintf("SELECT `id`, `tid`, `rid`, `info` FROM `tasks_%s` WHERE `stat`=1 order by `id` DESC LIMIT ?", ttype)
	return db.Query(sqlStr, num)
}


/* monitor >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> */
func InfoData(ttype, tid string) ([]map[string]string, error) {
	sqlStr := fmt.Sprintf("SELECT * FROM `tasks_%s` WHERE `tid`=?", ttype)
	rst, err := db.Query(sqlStr, tid)
	if err != nil {
		return nil, err
	}
	
	return rst, nil
}

func GetCount(ttype string) (ncout, icout, scout, ecout int64, err error) {
	var rst []map[string]string
	sqlStr := fmt.Sprintf("SELECT `stat`, count(0) cou FROM `tasks_%s` GROUP BY `stat`", ttype)
	rst, err = db.Query(sqlStr)
	if err != nil {
		return
	}

	for _, nv := range rst {
		ri, _ := strconv.ParseInt(nv["cou"], 10, 64)

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

func GetRecord(ttype string) (nr, ir, er []string, err error) {
	const NUM = 20
	var rst []map[string]string
	sqlStr := fmt.Sprintf("select tid,stat from(select @gid:=@cgid,@cgid:=t1.stat,if(@gid=@cgid,@rank:=@rank+1,@rank:=1) rank,t1.* FROM (select tid,stat from `tasks_%s` ORDER BY `id` DESC) t1,(select @gid:=1,@cgid:=1,@rank:=1) t2)t3 where t3.rank<=%d", ttype, NUM)
	rst, err = db.Query(sqlStr)
	if err != nil {
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
/* monitor <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< */
