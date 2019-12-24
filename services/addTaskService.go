package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/liuhengloveyou/easyTask/models"
	"github.com/liuhengloveyou/easyTask/rappers"
)

// 添加任务
func AddTaskService(uid int64, taskType string, data []byte, batch bool) (ids []int64, err error) {
	// 任务类型有吗？
	if _, err = rappers.NewRapper(taskType); err != nil {
		logger.Errorf("AddTaskService no rapper types: %s\n", taskType)
		return
	}

	if batch {
		return addTaskBatch(uid, taskType, data)
	} else {
		return addTaskOne(uid, taskType, data)
	}
}

// 添加任务
func addTaskOne(uid int64, taskType string, data []byte) (ids []int64, err error) {
	rapper, _ := rappers.NewRapper(taskType)

	info := rapper.NewTaskInfo()
	if err = json.Unmarshal([]byte(data), info); err != nil {
		logger.Errorf("AddTaskService Unmarshal body ERR: %v\n", err.Error())
		return
	}

	task := &models.Task{UID: uid}
	task.TaskType = taskType
	task.Stat = models.TaskStatusNew
	if info.GetRid() != "" {
		task.Rid = info.GetRid()
	}
	task.TaskInfo, _ = json.Marshal(info)

	logger.Debugf("AddTaskService model: %#v %#v\n", taskType, info)

	var id int64
	id, err = task.Insert()
	if err != nil {
		logger.Error("AddTaskService ERR: ", err.Error())
		return
	}

	logger.Infof("AddTaskService: %d %s\n", id, err)

	ids = append(ids, id)

	return
}

func addTaskBatch(uid int64, taskType string, data []byte) (ids []int64, err error) {
	rapper, _ := rappers.NewRapper(taskType)

	lines := strings.Split(string(data), "\r\n")
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}

		logger.Debugf("addTaskBatch line: %s\v", line)

		info := rapper.NewTaskInfo()
		if err = info.FromString(line); err != nil {
			logger.Errorf("addTaskBatch fmt ERR: %s %s\n", taskType, line)
			return nil, fmt.Errorf("数据格式错误")
		}

		task := &models.Task{UID: uid}
		task.TaskType = taskType
		task.Stat = models.TaskStatusNew
		if info.GetRid() != "" {
			task.Rid = info.GetRid()
		}
		task.TaskInfo, _ = json.Marshal(info)

		logger.Debugf("AddTaskService model: %#v %#v\n", taskType, info)

		var id int64
		id, err = task.Insert()
		if err != nil {
			logger.Error("AddTaskService ERR: ", err.Error())
			return nil, fmt.Errorf("数据入库错误")
		}

		ids = append(ids, id)
	}

	logger.Infof("AddTaskService OK: %#v %s\n", ids, err)

	return
}
