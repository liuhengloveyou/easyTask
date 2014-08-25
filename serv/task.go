package main

import (
	"sync"
	"time"
	"container/list"
	
	"github.com/golang/glog"
)

type TaskType struct {
	name     string
	rappers  map[string]*Rapper
	taskChan chan *TaskInfo

	backList *list.List // *TaskInfo
	inList   *list.List // *TaskInfo2DB
	backLock sync.Mutex
	inLock   sync.RWMutex

	once [2]sync.Once
}

type TaskInfo2DB struct {
	taskInfo *TaskInfo
	sign     byte // 'A' = insert; 'U' = update
	stat     int64
	msg      string
}

func NewTaskType() *TaskType {
	return new(TaskType).Init()
}

func (this *TaskType) Init() *TaskType {
	this.name = ""
	this.rappers = make(map[string]*Rapper)
	this.taskChan = make(chan *TaskInfo, int64(confJson["MaxTaskPerRapper"].(float64))*2)
	this.backList = list.New()
	this.inList = list.New()

	return this
}

func (this *TaskType) resetRapper(one *Rapper) {
	tasks := one.ReSet()
	one.Beat()
	
	this.backLock.Lock()
	for _, tn := range tasks {
		this.backList.PushFront(tn)
	}
	this.backLock.Unlock()

	return
}

func (this *TaskType) newTask(task *TaskInfo, stat int64) {
	this.once[1].Do(func() {
		go this.realUpTask()
	})

	t := &TaskInfo2DB{taskInfo: task, sign: 'A', stat: stat}

	this.inLock.Lock()
	this.inList.PushBack(t)
	this.inLock.Unlock()
}

func (this *TaskType) upTask(one *Rapper, stat int64, tid, msg string) {
	this.once[1].Do(func() {
		go this.realUpTask()
	})

	one.DelTask(tid)

	t := &TaskInfo2DB{taskInfo: &TaskInfo{Tid: tid}, stat: stat, msg: msg, sign: 'U'}
	this.inLock.Lock()
	this.inList.PushBack(t)
	this.inLock.Unlock()
}

func (this *TaskType) realUpTask() {
	for {
		one := this.inList.Front()
		if one == nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		t := one.Value.(*TaskInfo2DB)
		if t.sign == 'A' {
			_, err := newTask2DB(this.name, t.taskInfo.Tid, t.taskInfo.Rid, t.taskInfo.Info, t.stat)
			if err != nil {
				glog.Errorln(err)
		} else if t.stat == 2 {
				this.backLock.Lock()
				this.backList.PushBack(one)
				this.backLock.Unlock()
			}
		} else if t.sign == 'U' {
			_, err := upTask2DB(this.name, t.taskInfo.Tid, t.msg, t.stat)
			if err != nil {
				glog.Errorln(err)
			}
		}

		this.inLock.Lock()
		this.inList.Remove(one)
		this.inLock.Unlock()
	}
}

func (this *TaskType) distTask(one *Rapper, num int) []*TaskInfo {
	this.once[1].Do(func() {
		go this.realDistTask()
	})

	rst := make([]*TaskInfo, num)
	for i := 0; i < num; i++ {
		rst[i] = <-this.taskChan
		one.AddTask(rst[i])
	}

	return rst
}

func (this *TaskType) realDistTask() {
	for {
		one := this.backList.Front()
		if one != nil {
			this.taskChan <- one.Value.(*TaskInfo)
			this.backLock.Lock()
			this.backList.Remove(one) // 这里不能批量操作
			this.backLock.Unlock()
			continue
		}

		tasks, err := getTasks(this.name, int64(confJson["MaxTaskPerRapper"].(float64))*2)
		if err != nil {
			glog.Errorln(err)
		}

		if tasks == nil {
			this.taskChan <- nil // 写满阻塞
			continue
		}

		for _, tn := range tasks {
			this.taskChan <- &TaskInfo{Tid: tn["tid"], Rid: tn["rid"], Info: tn["info"]}
		}
	}
}

func (this *TaskType) BuffSize() (int64, int64) {
	return int64(this.inList.Len()), int64(this.backList.Len())
}
