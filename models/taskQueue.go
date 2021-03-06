package models

import (
	"fmt"
	"sync"
	"time"

	gocommon "github.com/liuhengloveyou/go-common"
)

var taskQueues map[string]*TaskQueue = make(map[string]*TaskQueue)

type TaskQueue struct {
	TaskType string

	taskChan chan Task
	once     sync.Once
}

func newTaskQueue(taskType string) *TaskQueue {
	o := &TaskQueue{
		TaskType: taskType,
	}
	o.taskChan = make(chan Task)

	return o
}

func GetTaskQueue(taskType string) *TaskQueue {
	if queue, ok := taskQueues[taskType]; ok {
		return queue
	} else {
		newQueue := newTaskQueue(taskType)
		taskQueues[taskType] = newQueue
		return newQueue
	}
}

func (this *TaskQueue) DistTask() Task {
	this.once.Do(func() {
		logger.Debug("go realDistTask...")
		go this.realDistTask()
	})

	task := <-this.taskChan

	return task
}

func (this *TaskQueue) realDistTask() {
	for {
		m := &Task{
			TaskType: this.TaskType,
			Stat:     TaskStatusNew,
		}

		tasks, err := m.Query(1, 10)
		if err != nil {
			logger.Error("query tasks ERR: ", this.TaskType, err.Error())
			time.Sleep(3 * time.Second) // 等会儿再查
		}

		if tasks == nil || len(tasks) < 1 {
			this.taskChan <- Task{ID: -1} // 放个空的
			continue
		}

		for _, task := range tasks {
			this.taskChan <- task
			logger.Debugf("task to queue: %v %v\n", task.ID, task.Rid)
			// 更新任务状态
			task.Stat = TaskStatusSend
			if rows, err := task.Update(); rows != 1 || err != nil {
				logger.Errorf("update task state to send ERR: %d %v \n", rows, err)
			}
			logger.Infof("update task to send OK: %d %v\n", task.ID, task.Rid)
		}
	}
}

func (this *TaskQueue) GetTaskFromServe(taskServeAddr, taskType, name string, num int) Task {
	this.once.Do(func() {
		go this.realGetTaskFromServe(taskServeAddr, taskType, name, num)
	})

	return <-this.taskChan
}

func (this *TaskQueue) realGetTaskFromServe(taskServeAddr, taskType, name string, num int) {
	urlStr := fmt.Sprintf("%s/api?type=%s&name=%s&num=%d", taskServeAddr, taskType, name, num)

	for {
		resp, body, err := gocommon.GetRequest(urlStr, map[string]string{"X-API": "/task/get"})
		if err != nil {
			logger.Error("realGetTaskFromServe ERR: ", urlStr, err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		if resp.StatusCode != 200 {
			logger.Error("realGetTaskFromServe StatusCode ERR: ", urlStr, resp.StatusCode)
			time.Sleep(1 * time.Second)
			continue
		}

		var tasks []Task
		if err = gocommon.UnmarshalHttpResponse(body, &tasks); err != nil {
			logger.Errorf("realGetTaskFromServe data ERR: ", urlStr, err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		if len(tasks) < 1 {
			logger.Errorf("realGetTaskFromServe data nil")
			this.taskChan <- Task{ID: -1} // 放个空的
			continue
		}

		for i := 0; i < len(tasks); i++ {
			if tasks[i].ID <= 0 || tasks[i].Rid == "" || tasks[i].TaskType != taskType {
				logger.Errorf("realGetTaskFromServe result ERR: %#v\n", tasks[i].ID, tasks[i].Rid)
				continue
			}

			this.taskChan <- tasks[i]
			logger.Infof("realGetTaskFromServe one: %v %v\n", tasks[i].ID, tasks[i].Rid)
		}

	}
}
