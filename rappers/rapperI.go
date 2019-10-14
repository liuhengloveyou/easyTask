package rappers

import "fmt"

// 系统里所有的rapper类型
var rappers = make(map[string]rapperType)

// rapper 需要实现的接口
type Rapper interface {
	Run() // 开始任务

	NewTaskInfo() interface{} // 返回一个任务详情对象
}

type rapperType func() Rapper

func RegisterRapper(name string, one rapperType) {
	if one == nil {
		panic("register rapper nil")
	}
	if _, dup := rappers[name]; dup {
		panic("register duplicate for " + name)
	}
	rappers[name] = one
}

func NewRapper(typeName string) (Rapper, error) {
	newFun, ok := rappers[typeName]
	if ok != true {
		return nil, fmt.Errorf("no rapper types " + typeName)
	}

	return newFun(), nil
}
