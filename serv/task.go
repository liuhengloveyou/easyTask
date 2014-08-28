package main

import (
	"container/list"
	"strconv"
	"sync"
	"time"

	"github.com/golang/glog"
)

type TaskType struct {
	name     string
	rappers  map[string]*Rapper
	taskChan chan *TaskInfo

	buff [2]*list.List // outbuff: *TaskInfo, inbuff: *TaskInfo2DB
	lock [2]sync.Mutex // outbuff lock, inbuff lock
	once [2]sync.Once  // outbuff once, inbuff once
}

type TaskInfo2DB struct {
	taskInfo *TaskInfo
	sign     byte // 'A' = insert; 'U' = update
	stat     int64
	rapper   string
	msg      string
}

func NewTaskType() *TaskType {
	return new(TaskType).Init()
}

func (this *TaskType) Init() *TaskType {
	this.name = ""
	this.rappers = make(map[string]*Rapper)
	this.taskChan = make(chan *TaskInfo, int64(confJson["MaxTaskPerRapper"].(float64))*2)
	this.buff[0] = list.New()
	this.buff[1] = list.New()

	return this
}

func (this *TaskType) resetRapper(one *Rapper) {
	tasks := one.ReSet()

	this.lock[0].Lock()
	defer this.lock[0].Unlock()

	for _, tn := range tasks {
		this.buff[0].PushFront(tn)
	}

	return
}

func (this *TaskType) newTask(task *TaskInfo, stat int64) {
	this.once[1].Do(func() {
		go this.realUpTask()
	})

	t := &TaskInfo2DB{taskInfo: task, sign: 'A', stat: stat}
	this.lock[1].Lock()
	this.buff[1].PushBack(t)
	this.lock[1].Unlock()
}

func (this *TaskType) upTask(one *Rapper, stat int64, tid, msg string) {
	this.once[1].Do(func() {
		go this.realUpTask()
	})

	t := &TaskInfo2DB{taskInfo: &TaskInfo{Tid: tid}, stat: stat, rapper: one.Name, msg: msg, sign: 'U'}
	this.lock[1].Lock()
	this.buff[1].PushBack(t)
	this.lock[1].Unlock()

	one.DelTask(tid)
}

func (this *TaskType) realUpTask() {
	for {
		one := this.buff[1].Front()
		if one == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		t := one.Value.(*TaskInfo2DB)
		if t.sign == 'A' {
			_, err := newTask2DB(this.name, t.taskInfo.Tid, t.taskInfo.Rid, t.taskInfo.Info, t.stat)
			if err != nil {
				glog.Errorln(err)
			} else if t.stat == 2 {
				this.lock[0].Lock()
				this.buff[0].PushBack(t.taskInfo) // 直接入分发队列
				this.lock[0].Unlock()
			}
		} else if t.sign == 'U' {
			_, err := upTask2DB(this.name, t.taskInfo.Tid, t.rapper, t.msg, t.stat)
			if err != nil {
				glog.Errorln(err)
			}
		}

		this.lock[1].Lock()
		this.buff[1].Remove(one)
		this.lock[1].Unlock()
	}
}

func (this *TaskType) distTask(one *Rapper, num int) []*TaskInfo {
	this.once[0].Do(func() {
		go this.realDistTask()
	})

	rst := make([]*TaskInfo, num)
	for i := 0; i < num; i++ {
		rst[i] = <-this.taskChan
		if rst[i] != nil {
			one.AddTask(rst[i])
		}
	}

	return rst
}

func (this *TaskType) realDistTask() {
	for {
		one := this.buff[0].Front()
		if one != nil {
			this.taskChan <- one.Value.(*TaskInfo)
			this.lock[0].Lock()
			this.buff[0].Remove(one) // 这里不能批量操作
			this.lock[0].Unlock()
			continue
		}

		tasks, err := getTasks(this.name, int64(confJson["MaxTaskPerRapper"].(float64))*2)
		if err != nil {
			glog.Errorln(err)
		}

		if tasks == nil || len(tasks) < 1 {
			this.taskChan <- nil // 写满阻塞
			continue
		}

		var ids, ide string
		for _, tn := range tasks {
			this.taskChan <- &TaskInfo{Tid: tn["tid"], Rid: tn["rid"], Info: tn["info"]}

			ide = tn["id"]
			if ids == "" {
				ids = tn["id"]
			}
		}

		// 更新任务状态
		idsi, _ := strconv.ParseInt(ids, 10, 64)
		idei, _ := strconv.ParseInt(ide, 10, 64)
		if _, err := upTaskStat2DB(this.name, idsi, idei); err != nil {
			glog.Errorln(err)
		}
	}
}

func (this *TaskType) BuffSize() (int64, int64) {
	return int64(this.buff[1].Len()), int64(this.buff[0].Len())
}
