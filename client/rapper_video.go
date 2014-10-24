package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

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
type videoTaskInfo struct {
	Flen     int32  // 文件长度
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

type rapperVideo struct {
	no          int64
	toTranscode chan *videoTaskInfo
	toUpload    chan *videoTaskInfo
	toUpdate    chan *videoTaskInfo
}

func init() {
	register("video", NewRapperVideo)
}

func NewRapperVideo() rapper {
	return &rapperVideo{
		no:          time.Now().UnixNano(),
		toTranscode: make(chan *videoTaskInfo, 5),
		toUpload:    make(chan *videoTaskInfo, 5),
		toUpdate:    make(chan *videoTaskInfo, 5)}
}

func (this *rapperVideo) run() {
	go this.transcode()
	go this.uploadMp4()
	go this.updateTask()

	for {
		// 取一个任务
		oneTask := oneTask()
		oneVideoTask := &videoTaskInfo{Tid: oneTask.Tid, Rid: oneTask.Rid}
		if oneVideoTask.Tid == "" || oneVideoTask.Rid == "" {
			glog.Errorln("taskERR: ", oneTask)
			continue
		}
		glog.Infof("deal oneTask: ", this.no, oneTask)

		taskInfo, err := base64.StdEncoding.DecodeString(oneTask.Info)
		if err != nil {
			oneVideoTask.err = fmt.Errorf("taskERR: %s", err.Error())
			goto END
		}

		err = json.Unmarshal(taskInfo, oneVideoTask)
		if err != nil {
			oneVideoTask.err = fmt.Errorf("downloadERR: %s", err.Error())
			goto END
		}

		// 下载文件
		glog.Infoln("download oneTask begin: ", this.no, oneVideoTask.toString())
		err = download(oneVideoTask.Url, confJson["tmpdir"].(string)+oneVideoTask.Tid)
		if err != nil {
			oneVideoTask.err = fmt.Errorf("downloadERR: %s", err.Error())
		}

	END: // 下载完成
		if oneVideoTask.err == nil {
			glog.Warningln("downloadOK: ", this.no, oneVideoTask.toString())
			this.toTranscode <- oneVideoTask
		} else {
			glog.Errorln("downloadERR:", this.no, oneVideoTask.toString())
			this.toUpdate <- oneVideoTask
		}
	}
}

func (this *rapperVideo) transcode() {
	const TRANFMT = "export LANGUAGE=en_US;ffmpeg -y -v error -i %s -movflags +faststart -vcodec libx264 -b:v 512k -acodec libvo_aacenc -ab 128k %s.mp4"
	const THUMFMT = "export LANGUAGE=en_US;ffmpeg -y -v error -i %s.mp4 -y -f  image2 -ss 8.0  -vframes 1 -s 120x80 %s.jpg"

	for {
		// 取一个任务
		oneVideoTask := <-this.toTranscode

		fn := confJson["tmpdir"].(string) + oneVideoTask.Tid
		transcodeCmd := fmt.Sprintf(TRANFMT, fn, fn)
		thumbnailCmd := fmt.Sprintf(THUMFMT, fn, fn)

		// 转码
		glog.Infoln("transcode: ", this.no, oneVideoTask.toString(), transcodeCmd)
		_, err := exec.Command("/bin/bash", "-c", transcodeCmd).Output()
		if err != nil {
			oneVideoTask.err = fmt.Errorf("transcodeERR: %s", err.Error())
			goto END
		}

		// 取缩略图
		glog.Infoln("thumbnail: ", this.no, oneVideoTask.toString(), thumbnailCmd)
		_, err = exec.Command("/bin/bash", "-c", thumbnailCmd).Output()
		if err != nil {
			glog.Errorln(err, thumbnailCmd)
		}

	END: // 转码完成
		if oneVideoTask.err == nil {
			glog.Warningln("transcodeOK: ", this.no, oneVideoTask.toString())
			this.toUpload <- oneVideoTask
		} else {
			glog.Errorln("transcodeERR: ", this.no, oneVideoTask.toString())
			this.toUpdate <- oneVideoTask
		}
	}
}

func (this *rapperVideo) uploadMp4() {
	for {
		var resp, jresp []byte

		// 取一个任务
		oneVideoTask := <-this.toUpload

		para := map[string]string{}
		// para from get to post
		urlOpara := strings.Split(oneVideoTask.Nurl, "?")
		if len(urlOpara) == 2 {
			for _, one := range strings.Split(urlOpara[1], "&") {
				kv := strings.Split(one, "=")
				if len(kv) == 2 {
					para[kv[0]] = kv[1]
				}
			}
		}

		// 上传新视频
		fi, err := os.Stat(confJson["tmpdir"].(string) + oneVideoTask.Tid + ".mp4")
		if err != nil {
			oneVideoTask.err = fmt.Errorf("uploadERR: %s", err.Error())
			goto END
		}
		para["flen"] = fmt.Sprintf("%d", fi.Size())

		glog.Infoln("uploadmp4: ", this.no, oneVideoTask.toString())
		resp, err = upload(oneVideoTask.Nurl, confJson["tmpdir"].(string)+oneVideoTask.Tid+".mp4", &para)
		if err != nil {
			oneVideoTask.err = fmt.Errorf("uploadERR: %s", err.Error())
			goto END
		}
		oneVideoTask.nfid = strings.Trim(string(resp), "\n ")
		if oneVideoTask.nfid == "" {
			oneVideoTask.err = fmt.Errorf("uploadERR: nfid nil")
			goto END
		}

		// 上传新视频缩略图
		fi, err = os.Stat(confJson["tmpdir"].(string) + oneVideoTask.Tid + ".jpg")
		if err != nil {
			glog.Errorln(oneVideoTask.toString(), err)
			goto END
		}
		para = map[string]string{"flen": fmt.Sprintf("%d", fi.Size())}

		glog.Infoln("uploadjpg: ", this.no, oneVideoTask.toString(), para)
		jresp, err = upload(oneVideoTask.Nurl, confJson["tmpdir"].(string)+oneVideoTask.Tid+".jpg", &para)
		if err != nil {
			glog.Errorln(oneVideoTask.toString(), err.Error())
			goto END
		}
		oneVideoTask.nimg = strings.Trim(string(jresp), "\n ")

	END: // 上传完成
		if oneVideoTask.err == nil {
			glog.Warningln("uploadOK: ", this.no, oneVideoTask.toString())
		} else {
			glog.Errorln("uploadERR: ", this.no, oneVideoTask.toString())
		}
		this.toUpdate <- oneVideoTask
	}
}

func (this *rapperVideo) updateTask() {
	for {
		// 取一个任务
		oneVideoTask := <-this.toUpdate

		// 回调
		para := map[string]string{"type": confJson["tasktype"].(string), "tid": oneVideoTask.Tid, "rid": oneVideoTask.Rid}
		if oneVideoTask.err != nil {
			para["msg"] = base64.StdEncoding.EncodeToString([]byte(oneVideoTask.err.Error()))
		}
		if oneVideoTask.nfid != "" {
			para["nfid"] = oneVideoTask.nfid
		}
		if oneVideoTask.nimg != "" {
			para["img"] = oneVideoTask.nimg
		}

		glog.Warningln("callbackTask: ", this.no, oneVideoTask.toString(), para)
		_, err := getRequest(oneVideoTask.Callback, &para)
		if err != nil {
			olderr := "NULL"
			if oneVideoTask.err != nil {
				olderr = oneVideoTask.err.Error()
			}
			oneVideoTask.err = fmt.Errorf("%s\ncallbackERR: %s", olderr, err.Error())
			glog.Errorln("updateTask callbackERR:", oneVideoTask.toString())
		}

		// 更新任务状态
		delete(para, "rid")
		delete(para, "nfid")
		delete(para, "img")
		para["stat"] = "1"
		para["name"] = confJson["rappername"].(string)
		if oneVideoTask.err != nil {
			para["stat"] = "-1"
			para["msg"] = base64.StdEncoding.EncodeToString([]byte(oneVideoTask.err.Error()))
		}

		glog.Infoln("updateTask: ", this.no, oneVideoTask.toString(), para)
		_, err = getRequest(confJson["taskServ"].(string)+"/uptask", &para)
		if err == nil {
			glog.Warningln("updateTaskOK: ", this.no, oneVideoTask.toString())
		} else {
			glog.Errorln("updateTaskERR:", this.no, oneVideoTask.toString(), para, err.Error())
		}

		glog.Flush()

		// 删除临时文件
		fn := confJson["tmpdir"].(string) + oneVideoTask.Tid
		os.Remove(fn)
		os.Remove(fn + ".mp4")
		os.Remove(fn + ".jpg")
	}
}

func (this *videoTaskInfo) toString() string {
	errstr := ""
	if this.err != nil {
		errstr = this.err.Error()
	}
	return fmt.Sprintf("{flen: %d; fid: %s; type: %s; url: %s; nurl: %s; callback: %s; tid: %s; rid: %s; nfid: %s; nimg: %s; err: %s}",
		this.Flen, this.Fid, this.Type, this.Url, this.Nurl, this.Callback, this.Tid, this.Rid, this.nfid, this.nimg, errstr)
}
