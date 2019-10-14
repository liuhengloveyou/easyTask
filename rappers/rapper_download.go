package rappers

import ()

// 任务描述信息
type DownloadTaskInfo struct {
	Tid  string `json:"tid"` // 任务ID
	Rid  string `json:"rid"` // 记录ID
	Type string `json:"ftype"` // 文件类型
	URL  string `json:"url"`  // 文件下载URL
	Flen int32  // 文件长度
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
	//uploadUrlConf := ""
	//callbackUrlConf := ""
	//if val, ok := confJson["uploadUrl"]; ok == true && val.(string) != "" {
	//	uploadUrlConf = val.(string)
	//}
	//if val, ok := confJson["callbackUrl"]; ok == true && val.(string) != "" {
	//	callbackUrlConf = val.(string)
	//}
	//
	//go this.transcode()
	//go this.uploadMp4()
	//go this.updateTask()
	//
	//for {
	//	// 取一个任务
	//	oneTask := oneTask()
	//	oneVideoTask := &videoTaskInfo{Tid: oneTask.Tid, Rid: oneTask.Rid}
	//	if oneVideoTask.Tid == "" || oneVideoTask.Rid == "" {
	//		glog.Errorln("taskERR: ", oneTask)
	//		continue
	//	}
	//	glog.Infof("deal oneTask: ", this.no, oneTask)
	//
	//	taskInfo, err := base64.StdEncoding.DecodeString(oneTask.Info)
	//	if err != nil {
	//		oneVideoTask.err = fmt.Errorf("taskERR: %s", err.Error())
	//		goto END
	//	}
	//
	//	err = json.Unmarshal(taskInfo, oneVideoTask)
	//	if err != nil {
	//		oneVideoTask.err = fmt.Errorf("downloadERR: %s", err.Error())
	//		goto END
	//	}
	//
	//	// 覆盖配置参数
	//	if uploadUrlConf != "" {
	//		oneVideoTask.Nurl = uploadUrlConf
	//	}
	//	if callbackUrlConf != "" {
	//		oneVideoTask.Callback = callbackUrlConf
	//	}
	//
	//	// 下载文件
	//	glog.Infoln("download oneTask begin: ", this.no, oneVideoTask.toString())
	//	err = download(oneVideoTask.Url, confJson["tmpdir"].(string)+oneVideoTask.Tid)
	//	if err != nil {
	//		oneVideoTask.err = fmt.Errorf("downloadERR: %s", err.Error())
	//	}
	//
	//END: // 下载完成
	//	if oneVideoTask.err == nil {
	//		glog.Warningln("downloadOK: ", this.no, oneVideoTask.toString())
	//		this.toTranscode <- oneVideoTask
	//	} else {
	//		glog.Errorln("downloadERR:", this.no, oneVideoTask.toString())
	//		this.toUpdate <- oneVideoTask
	//	}
	//}
}
