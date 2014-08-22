package main

import (
	"sync"
	"time"
	"strconv"
	"container/list"

	"github.com/golang/glog"
)

type TaskType struct {
	name string
	rappers map[string]*Rapper
	taskChan chan *TaskInfo
	
	backList *list.List // *TaskInfo
	inList *list.List // *TaskInfo2DB
	backLock sync.Mutex
	inLock sync.RWMutex
	
	once [2]sync.Once
}

type TaskInfo2DB struct {
	taskInfo *TaskInfo
	sign byte // 'A' = insert; 'U' = update
	idno int64
	stat int64
	msg string
}

func NewTaskType() *TaskType {
	return new(TaskType).Init()
}

func (this *TaskType) Init() *TaskType {
	this.name = ""
	this.rappers = make(map[string]*Rapper)
	this.taskChan = make(chan *TaskInfo, int64(confJson["MaxTaskPerRapper"].(float64)) * 2)
	this.backList = list.New()
	this.inList = list.New()
	
	return this
}

func (this *TaskType) resetRapper(one *Rapper) {
	one.Lock.Lock()
	defer one.Lock.Unlock()
	this.backLock.Lock()
	defer this.backLock.Unlock()
	
	for e := one.Tasks.Front(); e != nil; e = e.Next() {
		this.backList.PushFront(e.Value.(*TaskInfo))
	}
	one.Tasks.Init()
	
	return
}

func (this *TaskType) newTask(task *TaskInfo) {
	this.once[1].Do(func() {
		go this._upTask()
	})
	
	t := &TaskInfo2DB{taskInfo: task, sign:'A'}
	
	this.inLock.Lock()
	this.inList.PushBack(t)
	this.inLock.Unlock()
}

func (this *TaskType) upTask(one *Rapper, stat int64, tid, msg string) {
	this.once[1].Do(func() {
		go this._upTask()
	})
	
	t := &TaskInfo2DB{taskInfo: &TaskInfo{Tid: tid}, stat: stat, msg: msg, sign:'U'}
	
	this.inLock.Lock()
	this.inList.PushBack(t)
	this.inLock.Unlock()
}

func (this *TaskType) _upTask() {
	for {
		time.Sleep(200 * time.Millisecond)
		
		one := this.inList.Front()
		if one == nil {
			continue
		}

		t := one.Value.(*TaskInfo2DB)
		if t.sign == 'A' {
			no, err := newTask2DB(this.name, t.taskInfo.Tid, t.taskInfo.Rid, t.taskInfo.Info)
			if err != nil {
				glog.Errorln(err)
			}
			
			t.idno = no
			this.backLock.Lock()
			this.backList.PushBack(one)
			this.backLock.Unlock()
			
			this.inLock.Lock()
			this.inList.Remove(one)
			this.inLock.Unlock()
		} else if t.sign == 'U' {

		}
	}
}

func (this *TaskType) distTask(one *Rapper, num int) []*TaskInfo{
	this.once[1].Do(func() {
		go this._distTask()
	})
	
	rst := make([]*TaskInfo, num)
	for i := 0; i < num; i++ {
		rst[i] = <- this.taskChan
	}

	return rst
}

func (this *TaskType) _distTask() {
	var rno int64 = 0
	
	for {
		one := this.backList.Front()
		if one != nil {
			this.taskChan <- one.Value.(*TaskInfo)
			this.backLock.Lock()
			this.backList.Remove(one)
			this.backLock.Unlock()
			continue
		}

		tasks, err := getTasks(this.name, int64(confJson["MaxTaskPerRapper"].(float64)) * 2, rno)
		if err != nil {
			glog.Errorln(err)
		}

		if tasks == nil {
			this.taskChan <- nil
			continue
		}

		for _, tn := range tasks {
			this.taskChan <- &TaskInfo{Tid: tn["tid"], Rid: tn["rid"], Info: tn["info"]}
			no, err := strconv.Atoi(tn["id"])
			if err != nil {
				glog.Errorln(err)
				continue
			}
			rno = int64(no)
		}
	}
}
