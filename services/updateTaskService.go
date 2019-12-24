package services

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/liuhengloveyou/easyTask/models"
)

func UpdateTaskService(body []byte) error {
	task := &models.Task{}
	if e := json.Unmarshal(body, task); e != nil {
		logger.Errorf("UpdateTaskService body ERR: %v\n", e.Error())
		return e
	}

	return realUpdateTaskService(task)
}

func RedoTaskService(uid int64, body []byte) error {
	if uid <= 0 {
		return fmt.Errorf("请求错误")
	}

	task := &models.Task{}
	if e := json.Unmarshal(body, task); e != nil {
		logger.Errorf("UpdateTaskService body ERR: %v\n", e.Error())
		return e
	}
	task.UID = uid
	task.Stat = models.TaskStatusNew

	return realUpdateTaskService(task)
}

// 更新任务状态
func realUpdateTaskService(task *models.Task) error {
	logger.Debug("UpdateTaskService model: ", task)
	if task.ID <= 0 {
		return fmt.Errorf("UpdateTaskService no id")
	}
	istat, _ := strconv.ParseFloat(task.Stat, 64)
	if istat < 0 || istat > 100 {
		return fmt.Errorf("UpdateTaskService Stat val err")
	}
	if task.Rapper == "" {
		logger.Warnf("UpdateTaskService RapperName nil")
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
