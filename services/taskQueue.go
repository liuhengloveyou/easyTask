package services

import (
	"container/list"
	"sync"
	"time"

	"github.com/liuhengloveyou/easyTask/models"
)

var taskQueues map[string]*TaskQueue = make(map[string]*TaskQueue)

type TaskQueue struct {
	TaskType     string
	taskChan chan *models.Task

	// 避免每次有请求时候才查库，批量查出来排队在这里
	buff *list.List
	lock sync.Mutex
	once sync.Once
}

func NewTaskQueue(taskType string) *TaskQueue {
	o :=  &TaskQueue{}
	o.Init(taskType)
	return o
}

func (this *TaskQueue) Init(taskType string) {
	this.TaskType = taskType
	this.taskChan = make(chan *models.Task, 10)
	this.buff = list.New()
}


func (this *TaskQueue) DistTask(num int) []*models.Task {
	this.once.Do(func() {
		go this.realDistTask()
	})

	rst := make([]*models.Task, 0)
	for i := 0; i < num; i++ {
		taskOne := <-this.taskChan
		if taskOne != nil {
			rst = append(rst, taskOne)
		}
	}

	return rst
}

func (this *TaskQueue) realDistTask() {
	for {
		one := this.buff.Front()
		if one != nil {
			this.taskChan <- one.Value.(*models.Task)
			this.lock.Lock()
			this.buff.Remove(one) // 这里不能批量操作
			this.lock.Unlock()
			continue
		}

		m := &models.Task{
			TaskType: this.TaskType,
		}

		tasks, err := m.Query(10)
		if err != nil {
			logger.Error("query tasks ERR: ", this.TaskType, err.Error())
			time.Sleep(3 * time.Second) // 等会儿再查
		}

		if tasks == nil || len(tasks) < 1 {
			time.Sleep(time.Second) // 等会儿再查
			continue
		}

		for _, tn := range tasks {
			this.taskChan <- &tn
			logger.Debug("task to queue:", this.TaskType, tn)
		}
	}
}

func (this *TaskQueue) BuffSize() int {
	return this.buff.Len()
}