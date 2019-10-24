package services

import (
	"encoding/json"
	"fmt"
	"github.com/liuhengloveyou/easyTask/models"
)

// 更新任务状态
func UpdateTaskService(body []byte) error {
	task := &models.Task{}
	if e := json.Unmarshal(body, task); e != nil {
		logger.Errorf("UpdateTaskService body ERR: %v\n", e.Error())
		return e
	}

	logger.Debug("UpdateTaskService model: ", task)
	if task.ID <= 0 {
		return fmt.Errorf("UpdateTaskService no id")
	}
	if task.Stat <= models.TaskStatusNew || task.Stat >= models.TaskStatusEND {
		return fmt.Errorf("UpdateTaskService Stat err")
	}
	if task.Rapper == "" {
		return fmt.Errorf("UpdateTaskService RapperName nil")
	}

	rows, err := task.Update()
	if err != nil {
		logger.Error("UpdateTaskService ERR: %v\n", err.Error())
		return err
	}
	if rows != 1 {
		logger.Error("UpdateTaskService rows ERR: ", rows)
		return fmt.Errorf("更新错误 %d", rows)
	}

	logger.Infof("UpdateTaskService OK: %v %v %v %v", task.ID, task.Rid, task.Stat, task.UpdateTime)

	return nil
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
