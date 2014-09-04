package models

import (
	"strings"
	"time"
	
	. "easyTask/serv/common"
)

var TaskTypes map[string]*TaskType = make(map[string]*TaskType) // 系统中所有的任务类型

var db *DBmysql // 数据库操作封装

func init() {
	var err error

	// 连接数据库
	if db, err = new(DBmysql).Init(ConfJson["dbUrl"].(string)); err != nil {
		panic(err)
	}

	// 加载所有任务类型名
	if err := loadTaskType(); err != nil {
		panic(err)
	}
}

func loadTaskType() error {
	tables, err := ShowTables()
	if err != nil {
		return err
	}

	for i := 0; i < len(tables); i++ {
		if strings.HasPrefix(tables[i], "tasks_") {
			tname := tables[i][6:]
			TaskTypes[tname] = NewTaskType()
			TaskTypes[tname].Name = tname
		}
	}

	return nil
}

func GetRapper(ttype, name string) (taskType *TaskType, rapper *Rapper) {
	ok := false
	
	taskType, ok = TaskTypes[ttype]
	if ok == true {
		rapper, _ = taskType.rappers[name]
	}
	
	return
}


func RapperCleaner() {
	for {
		time.Sleep(1 * time.Second)
		
		for _, v := range TaskTypes {
			for k1, v1 := range v.rappers {
				if v1.Beat(false) < 0 {
					delete(v.rappers, k1) // 已经死了, 删之
					continue
				}
				if (time.Now().Unix() - v1.Beat(false)) > int64(ConfJson["RapperBeatOut"].(float64)) {
					v.ResetRapper(v1)
					v1.Kill()
				}
			}
		}
	}
}
