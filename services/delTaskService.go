package services

import (
	"fmt"

	"github.com/liuhengloveyou/easyTask/models"
)

// 删除自己的任务
func DeleteTaskService(taskid, uid int64) (err error) {
	if taskid <= 0 || uid <= 0 {
		return fmt.Errorf("请求错误")
	}

	task := &models.Task{UID: uid, ID: taskid}
	logger.Debugf("DeleteTaskService model: %#v\n", task)

	rst, err := task.Delete(nil)
	if err != nil {
		logger.Error("DeleteTaskService ERR: ", err.Error())
		return err
	}

	row, _ := rst.RowsAffected()
	if row != 1 {
		logger.Errorf("DeleteTaskService 0")
		return fmt.Errorf("删除失败")
	}

	logger.Infof("DeleteTaskService: %d %s\n", uid, taskid, err)

	return nil
}
