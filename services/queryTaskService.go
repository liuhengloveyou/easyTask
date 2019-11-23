package services

import (
	"database/sql"
	"github.com/liuhengloveyou/easyTask/models"
	"github.com/liuhengloveyou/easyTask/rappers"
)

// 取任务
func GetTaskService(taskType string, num int) (tasks []models.Task, err error) {
	// 任务类型有吗？
	if _, err = rappers.NewRapper(taskType); err != nil {
		logger.Errorf("QueryTaskService no rapper types: %s\n", taskType)
		return
	}

	taskQueue := models.GetTaskQueue(taskType)

	// 从队列里取
	tasks = make([]models.Task, 0)
	for i := 0; i < num; i++ {
		one := taskQueue.DistTask()
		if one.ID < 0 {
			break
		}
		tasks = append(tasks, one)
	}

	return tasks, nil
}

// 查任务信息
func QueryTaskService(uid int64, taskType string, pageNO, pageSize uint) (tasks []models.Task, err error) {
	model := models.Task{
		ID:       -1,
		UID:      sql.NullInt64{Valid:true, Int64:uid},
		TaskType: taskType,
		Stat: -1,
	}

	tasks, err = model.Query(pageNO, pageSize)

	return
}
