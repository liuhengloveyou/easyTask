package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/liuhengloveyou/easyTask/models"
	"github.com/liuhengloveyou/easyTask/rappers"
)

// 添加一个任务
func AddTaskService(body []byte) (id int64, err error) {
	task := &models.Task{}
	if e := json.Unmarshal(body, task); e != nil {
		logger.Errorf("AddTaskService Unmarshal body ERR: %v\n", e.Error())
		return -1, e
	}

	logger.Debug("AddTaskService model: ", task)

	// 任务类型有吗？
	var rapper rappers.Rapper
	if rapper, err = rappers.NewRapper(task.TaskType); err != nil {
		logger.Errorf("AddTaskService no rapper types: %s\n", task.TaskType)
		return
	}

	// info格式对吗?
	taskInfo := rapper.NewTaskInfo()
	if err = json.Unmarshal([]byte(task.TaskInfo), taskInfo); err != nil {
		logger.Errorf("AddTaskService info ERR: ", task)
		return
	}

	logger.Debugf("AddTaskService model: %#v\n%#v\n", task, taskInfo)

	now := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	task.Stat = 1
	task.AddTime = now
	task.UpdateTime = now

	id, err = task.Insert()
	if err != nil {
		logger.Error("AddTaskService ERR: ", err.Error())
		return
	}

	logger.Infof("AddTaskService: %d %s\n", id, err)

	return
}

// 取任务
func QueryTaskService(taskType string, num int) (tasks []models.Task, err error) {

	// 任务类型有吗？
	// var rapper rappers.Rapper
	if _, err = rappers.NewRapper(taskType); err != nil {
		logger.Errorf("QueryTaskService no rapper types: %s\n", taskType)
		return
	}

	model := models.Task{
		TaskType: taskType,
	}

	tasks, err = model.Query(num)
	if err != nil {
		logger.Errorf("QueryTaskService db ERR:", err.Error())
		return nil, fmt.Errorf("服务错误")
	}

	return tasks, nil
}

// 更新任务状态
func UpdateTaskService() (tasks []models.Task, err error) {

	return nil, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////
//func GetRapper(ttype, name string) (taskType *TaskType, rapper *Rapper) {
//	ok := false
//
//	taskType, ok = taskTypes[ttype]
//	if ok == true {
//		rapper, _ = taskType.rappers[name]
//	}
//
//	return
//}
//
//// 清理太久不活动的工兵
//func RapperCleaner() {
//	for {
//		time.Sleep(1 * time.Second)
//
//		for _, v := range taskTypes {
//			for k1, v1 := range v.rappers {
//				if v1.Beat(false) < 0 {
//					delete(v.rappers, k1) // 已经死了, 删之
//					continue
//				}
//				if (time.Now().Unix() - v1.Beat(false)) > 100 {
//					v.ResetRapper(v1)
//					v1.Kill()
//				}
//			}
//		}
//	}
//}