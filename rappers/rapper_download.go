package rappers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/liuhengloveyou/easyTask/common"
	"github.com/liuhengloveyou/easyTask/models"

	gocommon "github.com/liuhengloveyou/go-common"
)

// 任务描述信息
type DownloadTaskInfo struct {
	Rid  string `json:"rid"`   // 记录ID
	Type string `json:"ftype"` // 文件类型
	URL  string `json:"url"`   // 文件下载URL
}


func (p *DownloadTaskInfo) GetRid() string {
	return p.Rid
}

func (p *DownloadTaskInfo) FromString(raw string) error {
	return nil
}

type DownloadRapper struct {
}

func NewDownloadRapper() Rapper {
	return &DownloadRapper{}
}

func (p *DownloadRapper) NewTaskInfo() TaskInfoI {
	return &DownloadTaskInfo{}
}

func (this *DownloadRapper) Run() {
	taskQueue := models.GetTaskQueue("download")

	for {
		// 取一个任务
		oneTask := taskQueue.GetTaskFromServe(common.ClientConfig.TaskServeAddr,
			common.ClientConfig.TaskType,
			common.ClientConfig.Name,
			1)

		taskInfo := &DownloadTaskInfo{}
		if err := json.Unmarshal(oneTask.TaskInfo, taskInfo); err != nil {
			common.Logger.Sugar().Errorf("task info ERR: ", oneTask, err.Error())
			time.Sleep(time.Second)
			continue
		}
		common.Logger.Sugar().Debugf("one newtask: %#v %v %v\n", oneTask.ID, oneTask.Rid, taskInfo.URL)

		domain, port, proxyAuth, err := common.GetProxyURL()
		if err != nil {
			common.Logger.Sugar().Errorf("PROXY ERR: ", err.Error())
		}

		// 下载文件
		downer := gocommon.Downloader{
			Headers: map[string]string{
				"X-Requested-With": "XMLHttpRequest",
				"Referer":          "http://music.migu.cn/v3/music/playlist",
				"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
			},
			URL:     taskInfo.URL,
			DstPath: "./dist/" + oneTask.Rid + taskInfo.Type,
		}
		if domain != "" && port != "" {
			downer.ProxyUrl = "http://"+domain+":" + port
			downer.ProxyAuth = proxyAuth
		}
		common.Logger.Sugar().Infof("downing %v %v\n", downer.ProxyUrl, downer.URL)
		resp, err := downer.Download(context.Background())
		if err != nil || resp.StatusCode != http.StatusOK {
			oneTask.Stat = models.TaskStatusERR
			oneTask.Remark = err.Error()
		} else {
			oneTask.Stat = models.TaskStatusOK
			oneTask.Remark = oneTask.Rid + taskInfo.Type
		}

		common.Logger.Sugar().Debugf("close proxy %v\n", domain + ":" + port)
		common.CloseProxyURL(port)

		oneTask.Rapper = common.ClientConfig.Name
		if err = UpdateTaskToServe(oneTask); err != nil {
			common.Logger.Sugar().Errorf("update task ERR: %v %v %v %v\n\n", oneTask.ID, oneTask.Rid, oneTask.Stat, err.Error())
		}

		common.Logger.Sugar().Infof("task end: %v %v %v %v %v\n\n", oneTask.ID, oneTask.Rid, oneTask.Stat, err, taskInfo.URL)
		//time.Sleep(1 * time.Second)
	}
}
