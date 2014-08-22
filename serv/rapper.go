package main

import (
	"sync"
	"time"
	"container/list"
)

type Rapper struct {
	Name string
	Tasks *list.List // *TaskInfo
	Lock sync.Mutex
	beat int64
}

func NewRapper() *Rapper {
	return new(Rapper).Init()
}

func (this *Rapper) Init() *Rapper {
	this.Name = ""
	this.Tasks = list.New()
	this.Lock = sync.Mutex{}
	this.beat = time.Now().Unix()
	return this
}

func (this *Rapper) Beat() int64 {
	if this.beat > 0 {
		this.beat = time.Now().Unix()
	}
	
	return this.beat
}

func (this *Rapper) Kill(){
	this.beat = -1
}

func (this *Rapper) TaskSize() int64 {
	return int64(this.Tasks.Len())
}

func (this *Rapper) AddTask(taskID *TaskInfo) {
	this.Lock.Lock()
	this.Tasks.PushBack(taskID)
	this.Lock.Unlock()
}

func (this *Rapper) GetTaskOne() string {
	e := this.Tasks.Back()
	if nil != e {
		return e.Value.(string)
	}

	return ""
}

func (this *Rapper) HasTask(taskID string) interface{} {
	if "" == taskID {
		return nil
	}

	for e := this.Tasks.Front(); e != nil; e = e.Next() {
		if e.Value == taskID {
			return e
		}
	}
	
	return nil
}
