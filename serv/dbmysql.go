package main

/*
CREATE DATABASE `taskManager` DEFAULT CHARACTER SET utf8 COLLATE utf8_bin;

CREATE TABLE `tasks_demo` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `tid` varchar(33) NOT NULL,
  `rid` varchar(32) NOT NULL,
  `info` varchar(1024) NOT NULL,
  `stat` int(11) NOT NULL DEFAULT '0',
  `addTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `getTime` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  `overTime` timestamp NULL DEFAULT '0000-00-00 00:00:00',
  `rapper` varchar(256) DEFAULT NULL,
  `client` varchar(256) DEFAULT NULL,
  `remark` text,
  PRIMARY KEY (`id`),
  UNIQUE KEY `inx_tid` (`tid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


stat:
1 = 新任务
2 = 正在处理
3 = 处理成功
-1 = 处理出错

*/

import (
	"fmt"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const TABLESQL = "CREATE TABLE `tasks_%s` (`id` int(11) NOT NULL AUTO_INCREMENT,`tid` varchar(33) NOT NULL,`rid` varchar(32) NOT NULL,`info` varchar(1024) NOT NULL,`stat` int(11) NOT NULL DEFAULT '0',`addTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,`getTime` timestamp NULL DEFAULT '0000-00-00 00:00:00',`overTime` timestamp NULL DEFAULT '0000-00-00 00:00:00',`rapper` varchar(256) DEFAULT NULL,`client` varchar(256) DEFAULT NULL,`remark` text,PRIMARY KEY (`id`),UNIQUE KEY `inx_tid` (`tid`)) ENGINE=InnoDB DEFAULT CHARSET=utf8"

var mysqlConn *sql.DB

func dbInit() error {
	var err error
	mysqlConn, err = sql.Open("mysql", confJson["mysqlUrl"].(string))
	if err != nil {
		return err
	}

	return nil
}

func createDB(name string) error {
	sqlStr := fmt.Sprintf(TABLESQL, name)
	_, err := mysqlConn.Exec(sqlStr)
	if err != nil {
		return err
	}

	return nil
}

func showTables() ([]string, error) {
	tables, err := doQuery("show tables")
	if err != nil {
		return nil, err
	}
	
	rst := make([]string, len(tables))
	for n, tn := range tables {
		rst[n] = tn["Tables_in_taskManager"]
	}

	return rst, nil
}

func newTask2DB(ttype, tid, rid, info string) (int64, error) {
	sqlStr := fmt.Sprintf("INSERT INTO `tasks_%s`(`tid`, `rid`, `info`, `stat`) VALUES(?,?,?, 1)", ttype);
	return doInsert(sqlStr, tid, rid, info)
}

func upTask2DB(ttype, tid, msg string, stat int64) (int64, error) {
	sqlStr := fmt.Sprintf("UPDATE `tasks_%s` SET `stat`=?, ``, `info`, `stat`) VALUES(?,?,?, 1)", ttype);
	return doInsert(sqlStr, tid, rid, info)
}

func getTasks(ttype string, num, fid int64) ([]map[string]string, error) {
	sqlStr := fmt.Sprintf("SELECT id, tid, rid info FROM `tasks_%s` WHERE `stat`=1 AND `id` > ? order by `addTime` limit ?", ttype)
	return doQuery(sqlStr, fid, num)
}

func doQuery(sqlStr string, args ...interface{}) ([]map[string]string, error) {
	stmt, err := mysqlConn.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var ret []map[string]string
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		
		tmap := make(map[string]string, len(columns))
		for i, col := range values {
			if col == nil {
				tmap[columns[i]] = ""
			} else {
				tmap[columns[i]] = string(col)
			}
		}
		ret = append(ret, tmap)
	}
	
	return ret, nil
}


func doInsert(sqlStr string, args ...interface{}) (int64, error) {
	stmt, err := mysqlConn.Prepare(sqlStr)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	rst, err := stmt.Exec(args...)
	if err != nil {
		return -1, err
	}

	return rst.LastInsertId()
}

func doUpdate(sqlStr string, args ...interface{}) (int64, error) {
	stmt, err := mysqlConn.Prepare(sqlStr)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	rst, err := stmt.Exec(args...)
	if err != nil {
		return -1, err
	}

	return rst.RowsAffected()
}
