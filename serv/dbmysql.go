package main

/*
CREATE DATABASE `taskManager` DEFAULT CHARACTER SET utf8 COLLATE utf8_bin;

CREATE TABLE `taskManager`.`tasks-demo` (
  `tid` CHAR(33) NOT NULL,
  `rid` VARCHAR(32) NOT NULL,
  `info` VARCHAR(1024) NOT NULL,
  `stat` INT NOT NULL DEFAULT 0,
  `addTime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `getTime` TIMESTAMP NULL DEFAULT 0,
  `overTime` TIMESTAMP NULL DEFAULT 0,
  `rapper` VARCHAR(256) NULL,
  `remark` TEXT NULL,
  PRIMARY KEY (`tid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

stat:
1 = 新任务
2 = 正在处理
3 = 成功处理完成
-1 = 处理出错

*/

import (
	"fmt"
	"database/sql"
	
	_ "github.com/go-sql-driver/mysql"
)

const TABLESQL = "CREATE TABLE `tasks-%s`(`tid` CHAR(33) NOT NULL,`rid` VARCHAR(32) NOT NULL,`info` VARCHAR(1024) NOT NULL,`stat` INT NOT NULL DEFAULT 0,`addTime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,`getTime` TIMESTAMP NULL DEFAULT 0,`overTime` TIMESTAMP NULL DEFAULT 0,`rapper` VARCHAR(256) NULL,`remark` TEXT NULL,PRIMARY KEY (`tid`)) ENGINE=InnoDB DEFAULT CHARSET=utf8"

var mysqlConn *sql.DB

func dbInit() {
	var err error
	mysqlConn, err = sql.Open("mysql", confJson["mysqlUrl"].(string))
	if err != nil {
		panic(err)
	}
}

func createDB(name string) error {
	sqlStr := fmt.Sprintf(TABLESQL, name)
	_, err := mysqlConn.Exec(sqlStr)
	if err != nil {
		return err
	}
	
	return nil
}

func newTask(ttype, tid, rid, info string) error {
	sqlStr := "INSERT INTO `tasks-" + ttype + "`(`taskid`, `rid`, `info`, `stat`) VALUES('" + tid + "','" + rid + "','" + info + "', 1)";
	_, err := mysqlConn.Exec(sqlStr)
	if err != nil{
		return err
	}

	return nil
}

