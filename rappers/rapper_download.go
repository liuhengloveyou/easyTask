package rappers

import (
	"encoding/json"
	"time"
	"net/http"
	"context"

	"github.com/liuhengloveyou/easyTask/common"
	"github.com/liuhengloveyou/easyTask/models"

	gocommon "github.com/liuhengloveyou/go-common"
)

const TASK_TYPE = "download"

// 任务描述信息
type DownloadTaskInfo struct {
	Rid  string `json:"rid"`   // 记录ID
	Type string `json:"ftype"` // 文件类型
	URL  string `json:"url"`   // 文件下载URL
}

type DownloadRapper struct {
}

func init() {
	RegisterRapper("download", NewDownloadRapper)
}

func NewDownloadRapper() Rapper {
	return &DownloadRapper{}
}

func (p *DownloadRapper) NewTaskInfo() interface{} {
	return &DownloadTaskInfo{}
}

func (this *DownloadRapper) Run() {
	taskQueue := models.GetTaskQueue(TASK_TYPE)

	for {
		// 取一个任务
		oneTask := taskQueue.GetTaskFromServe(common.ClientConfig.TaskServeAddr, common.ClientConfig.TaskType, 1)

		taskInfo := &DownloadTaskInfo{}
		if err := json.Unmarshal(oneTask.TaskInfo, taskInfo); err != nil {
			common.Logger.Sugar().Errorf("task info ERR: ", oneTask, err.Error())
			time.Sleep(time.Second)
			continue
		}
		common.Logger.Sugar().Debugf("one newtask: %#v %v %v\n", oneTask.ID, oneTask.Rid, taskInfo.URL)

		// 下载文件
		downer := gocommon.Downloader{
			URL: taskInfo.URL,
			DstPath: "./dist/" + oneTask.Rid + taskInfo.Type,
		}
		common.Logger.Sugar().Infof("downing %v\n", downer.URL)
		resp, err := downer.Download(context.Background())
		if err != nil || resp.StatusCode != http.StatusOK {
			oneTask.Stat = models.TaskStatusERR
			oneTask.Remark = err.Error() + resp.Status
		}else {
			oneTask.Stat = models.TaskStatusOK
			oneTask.Remark = oneTask.Rid + taskInfo.Type
		}

		err = UpdateTaskToServe(oneTask)
		if err != nil {
			panic(err)
		}

		common.Logger.Sugar().Infof("task end: %v %v %v %v\n\n", oneTask.ID, oneTask.Rid, oneTask.Stat, taskInfo.URL)
	}
}
