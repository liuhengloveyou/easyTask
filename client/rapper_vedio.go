package main

import (
	"os"
	"fmt"
	"time"
	"os/exec"
	"encoding/json"
	"encoding/base64"

	"github.com/golang/glog"
)

/*
{
 "fid":"900499",
 "flen":15070734,
 "type":"wmv",
 "url":"http://www.upload.cn/client14081804173567259325.wmv",
 "nurl":"http://192.168.1.2:8081/ffaceup",
 "callback":"http://www.space.cn/opus/updateNfid.html"
}
*/
type vedioTaskInfo struct {
	Flen     int32 // 文件长度
	Fid      string // 文件ID
	Type     string // 文件类型
	Url      string // 文件下载URL
	Nurl     string // 文件长度
	Callback string // 任务回调
	Tid      string // 任务ID
	Rid      string // 记录ID
	nfid     string // 转码处理后的文件ID
	nimg     string // 缩略图 
	err      error  // 处理过程中的错误信息
}

type rapperVedio struct {
	no int64
	toTranscode chan *vedioTaskInfo
	toUpload chan *vedioTaskInfo
	toUpdate chan *vedioTaskInfo
}

func init() {
	register("vedio", NewRapperVedio)
}

func NewRapperVedio() rapper {
	return &rapperVedio{
		no: time.Now().UnixNano(),
		toTranscode: make(chan *vedioTaskInfo, 3),
		toUpload: make(chan *vedioTaskInfo, 3),
		toUpdate: make(chan *vedioTaskInfo, 3)}
}

func (this *rapperVedio) run() {
	go this.transcode()
	go this.uploadMp4()
	go this.updateTask()
	
	for {
		oneVedioTask := &vedioTaskInfo{}

		// 取一个任务
		oneTask := oneTask()
		glog.Infoln("download oneTask: ", this.no, oneTask)
		oneVedioTask.Tid = oneTask.Tid
		oneVedioTask.Rid = oneTask.Rid
		if oneVedioTask.Tid == "" || oneVedioTask.Rid == "" {
			glog.Errorln("taskERR: ", oneTask)
			continue
		}

		taskInfo, err := base64.StdEncoding.DecodeString(oneTask.Info)
		if err != nil {
			oneVedioTask.err = err
			glog.Errorln(err, oneTask)
			goto END
		}

		err = json.Unmarshal(taskInfo, oneVedioTask)
		if err != nil {
			oneVedioTask.err = err
			glog.Errorln(err, oneVedioTask)
			goto END
		}
		glog.Infoln("download oneTask OK: ", oneVedioTask)

		// 下载文件
		err = download(oneVedioTask.Url, confJson["tmpdir"].(string) + oneVedioTask.Tid)
		if err != nil {
			oneVedioTask.err = err
			glog.Errorln(oneVedioTask, err)
			goto END
		}

	END: // 下载完成
		if oneVedioTask.err == nil {
			glog.Infoln("downloadOK: ", oneVedioTask)
			this.toTranscode <- oneVedioTask
		} else {
			glog.Infoln("downloadERR: ", oneVedioTask, oneVedioTask.err)
			this.toUpdate <- oneVedioTask
		}

		// time.Sleep(3 * time.Second)
	}
}

func (this *rapperVedio) transcode() {
	const TRANFMT = "export LANGUAGE=en_US;ffmpeg -v error -i %s -movflags +faststart -vcodec libx264 -b:v 512k -acodec libvo_aacenc -ab 128k %s.mp4"
	const THUMFMT = "export LANGUAGE=en_US;ffmpeg -v error -i %s.mp4 -y -f  image2 -ss 8.0  -vframes 1 -s 120x80 %s.jpg"

	for {
		// 取一个任务
		oneVedioTask := <- this.toTranscode
		fn := confJson["tmpdir"].(string) + oneVedioTask.Tid

		transcodeCmd := fmt.Sprintf(TRANFMT, fn, fn)
		thumbnailCmd := fmt.Sprintf(THUMFMT, fn, fn)

		// 转码
		glog.Infoln("transcode:", transcodeCmd)
		_, err := exec.Command(transcodeCmd).Output()
		if err != nil {
			oneVedioTask.err = err
			glog.Errorln(err, transcodeCmd)
			goto END
		}

		// 取缩略图
		glog.Infoln("thumbnail:", thumbnailCmd)
		_, err = exec.Command(thumbnailCmd).Output()
		if err != nil {
			glog.Errorln(err, thumbnailCmd)
		}

	END: // 转码完成
		if oneVedioTask.err == nil {
			glog.Infoln("transcodeOK: ", oneVedioTask)
			this.toUpdate <- oneVedioTask
		} else {
			glog.Infoln("transcodeERR: ", oneVedioTask, oneVedioTask.err)
			this.toUpdate <- oneVedioTask
		}
	}

}

func (this *rapperVedio) uploadMp4() {
	for {
		var (
			err         error
			fi          os.FileInfo
			jfp         *os.File
			para        *map[string]string
			resp, jresp []byte
		)
		
		// 取一个任务
		oneVedioTask := <- this.toUpload
		fn := confJson["tmpdir"].(string) + oneVedioTask.Tid

		// 上传新视频
		fp, err := os.Open(fn+".mp4")
		if err != nil {
			oneVedioTask.err = err
			glog.Errorln(oneVedioTask, err)
			goto END

		}
		fi, err = fp.Stat()
		if err != nil {
			oneVedioTask.err = err
			glog.Errorln(oneVedioTask, err)
			goto END
		}
		para = &map[string]string{"flen": fmt.Sprintf("%d", fi.Size())}

		resp, err = upload(oneVedioTask.Nurl, fn+".mp4", para)
		if err != nil {
			oneVedioTask.err = err
			glog.Errorln(err, oneVedioTask)
			goto END
		}
		oneVedioTask.nfid = string(resp)
		
		// 上传新视频缩略图
		jfp, err = os.Open(fn+".jpg")
		if err != nil {
			glog.Errorln(oneVedioTask, err)
			goto END

		}
		fi, err = jfp.Stat()
		if err != nil {
			glog.Errorln(oneVedioTask, err)
			goto END
		}
		(*para)["flen"] = fmt.Sprintf("%d", fi.Size())

		jresp, err = upload(oneVedioTask.Nurl, fn+".jpg", para)
		if err != nil {
			glog.Errorln(err, oneVedioTask)
			goto END
		}
		oneVedioTask.nimg = string(jresp)
		
	END: // 上传完成
		if fp != nil {
			fp.Close()
		}
		if jfp != nil {
			jfp.Close()
		}
		if oneVedioTask.err == nil {
			glog.Infoln("uploadOK: ", oneVedioTask)
		} else {
			glog.Infoln("uploadERR: ", oneVedioTask, oneVedioTask.err)
		}
		this.toUpdate <- oneVedioTask
	}
}

func (this *rapperVedio) updateTask() {
	for {
		// 取一个任务
		oneVedioTask := <- this.toUpdate
		glog.Infoln("updateTask:", oneVedioTask)

		// 回调
		para := &map[string]string{"type": confJson["tasktype"].(string), "tid": oneVedioTask.Tid, "rid": oneVedioTask.Rid, "nfid": oneVedioTask.nfid}
		if oneVedioTask.nimg != "" {
			(*para)["img"] = oneVedioTask.nimg
		} 
		if oneVedioTask.err != nil {
			(*para)["msg"] = base64.StdEncoding.EncodeToString([]byte(oneVedioTask.err.Error()))
		}
		_, err := getRequest(oneVedioTask.Callback, para)
		if err != nil {
			oneVedioTask.err = err
			glog.Errorln("updateTask callbackERR:", oneVedioTask, err)
		}
		
		// 更新任务状态
		delete(*para, "tid")
		delete(*para, "rid")
		delete(*para, "nfid")
		delete(*para, "img")
		(*para)["me"] = confJson["rappername"].(string)
		(*para)["stat"] = "0"
		if oneVedioTask.err != nil {
			(*para)["stat"] = "-1"
			(*para)["msg"] = base64.StdEncoding.EncodeToString([]byte(oneVedioTask.err.Error()))
		}
		_, err = getRequest(confJson["taskServ"].(string) + "/uptask", para)
		if err != nil {
			oneVedioTask.err = err
			glog.Errorln("updateTask updateERR:", oneVedioTask, err)
		}

		// 删除临时文件
		fn := confJson["tmpdir"].(string) + oneVedioTask.Tid
		os.Remove(fn)
		os.Remove(fn+".mp4")
		os.Remove(fn+".jpg")
	}	
}
