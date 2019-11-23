package rappers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/liuhengloveyou/easyTask/common"
	"github.com/liuhengloveyou/easyTask/models"

	gocommon "github.com/liuhengloveyou/go-common"
)

// 系统里所有的rapper类型
var rappers = make(map[string]rapperType)

// taskInfo 需要实现的接口
type TaskInfoI interface {
	GetRid() string // 返回RID
	FromString(raw string) error // 从字符串生成
}

// rapper 需要实现的接口
type Rapper interface {
	Run()                   // 开始任务
	NewTaskInfo() TaskInfoI // 返回一个任务详情对象
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

func init() {
	RegisterRapper("download", NewDownloadRapper)
	RegisterRapper("daikan", NewDaikanRapper)
}

//////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////

// 客户端更新任务状态
func UpdateTaskToServe(task models.Task) error {
	if task.ID <= 0 || (task.Stat < models.TaskStatusNew || task.Stat > models.TaskStatusERR) {
		return fmt.Errorf("UpdateTaskToServe param ERR")
	}

	body, _ := json.Marshal(task)

	resp, _, err := gocommon.PostRequest(common.ClientConfig.TaskServeAddr+"/api", map[string]string{"X-API": "/task/update"}, nil, body)
	if err != nil {
		return fmt.Errorf("updateTaskToServe ERR: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("updateTaskToServe status ERR: %d", resp.StatusCode)
	}

	return nil
}
