package models

import (
	"sync"
	"time"
)

type Rapper struct {
	Name  string
	tasks map[string]*TaskInfo
	lock  sync.Mutex
	beat  int64
}

func (this *Rapper) Init() *Rapper {
	this.Name = ""
	this.tasks = make(map[string]*TaskInfo)
	this.lock = sync.Mutex{}
	this.beat = time.Now().Unix()
	return this
}

func NewRapper(name string) *Rapper {
	one := new(Rapper).Init()
	one.Name = name
	return one
}

func (this *Rapper) Beat(b bool) int64 {
	if b == true {
		this.beat = time.Now().Unix()
	}

	return this.beat
}

func (this *Rapper) Kill() {
	this.beat = -1
}

func (this *Rapper) TaskSize() int {
	return len(this.tasks)
}

func (this *Rapper) AddTask(task *TaskInfo) {
	this.lock.Lock()
	this.tasks[task.Tid] = task
	this.lock.Unlock()
}

func (this *Rapper) DelTask(tid string) {
	this.lock.Lock()
	delete(this.tasks, tid)
	this.lock.Unlock()
}

func (this *Rapper) ReSet() []*TaskInfo {
	this.lock.Lock()
	taskLen := len(this.tasks)
	tasks := make([]*TaskInfo, taskLen)
	for _, v := range this.tasks {
		tasks[taskLen-1] = v
		taskLen = taskLen - 1
	}

	this.tasks = make(map[string]*TaskInfo)
	this.lock.Unlock()

	return tasks
}
