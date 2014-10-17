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
		toTranscode: make(chan *videoTaskInfo, 3),
		toUpload:    make(chan *videoTaskInfo, 3),
		toUpdate:    make(chan *videoTaskInfo, 3)}
}

func (this *rapperVideo) run() {
	go this.transcode()
	go this.uploadMp4()
	go this.updateTask()

	for {
		oneVideoTask := &videoTaskInfo{}

		// 取一个任务
		oneTask := oneTask()
		glog.Infoln("download oneTask: ", this.no, oneTask.Tid, oneTask)
		oneVideoTask.Tid = oneTask.Tid
		oneVideoTask.Rid = oneTask.Rid
		if oneVideoTask.Tid == "" || oneVideoTask.Rid == "" {
			glog.Errorln("taskERR: ", oneTask)
			continue
		}

		taskInfo, err := base64.StdEncoding.DecodeString(oneTask.Info)
		if err != nil {
			oneVideoTask.err = fmt.Errorf("downloadERR: %s", err.Error())
			goto END
		}

		err = json.Unmarshal(taskInfo, oneVideoTask)
		if err != nil {
			oneVideoTask.err = fmt.Errorf("downloadERR: %s", err.Error())
			goto END
		}

		// 下载文件
		glog.Infoln("download oneTask begin: ", oneVideoTask)
		err = download(oneVideoTask.Url, confJson["tmpdir"].(string)+oneVideoTask.Tid)
		if err != nil {
			oneVideoTask.err = fmt.Errorf("downloadERR: %s", err.Error())
		}

	END: // 下载完成
		if oneVideoTask.err == nil {
			glog.Infoln("downloadOK: ", oneVideoTask)
			this.toTranscode <- oneVideoTask
		} else {
			glog.Errorln("downloadERR:", oneVideoTask, oneVideoTask.err.Error())
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
		glog.Infoln("transcode:", oneVideoTask, transcodeCmd)
		_, err := exec.Command("/bin/bash", "-c", transcodeCmd).Output()
		if err != nil {
			oneVideoTask.err = fmt.Errorf("transcodeERR: %s", err.Error())
			goto END
		}

		// 取缩略图
		glog.Infoln("thumbnail:", oneVideoTask, thumbnailCmd)
		_, err = exec.Command("/bin/bash", "-c", thumbnailCmd).Output()
		if err != nil {
			glog.Errorln(err, thumbnailCmd)
		}

	END: // 转码完成
		if oneVideoTask.err == nil {
			glog.Infoln("transcodeOK: ", oneVideoTask)
			this.toUpload <- oneVideoTask
		} else {
			glog.Errorf("transcodeERR: %v. '%s'", oneVideoTask, oneVideoTask.err.Error())
			this.toUpdate <- oneVideoTask
		}
	}
}

func (this *rapperVideo) uploadMp4() {
	for {
		var (
			err         error
			fi          os.FileInfo
			jfp         *os.File
			para        *map[string]string
			resp, jresp []byte

			urlOpara, kv []string
			one          string
		)

		// 取一个任务
		oneVideoTask := <-this.toUpload
		fn := confJson["tmpdir"].(string) + oneVideoTask.Tid

		// 上传新视频
		fp, err := os.Open(fn + ".mp4")
		if err != nil {
			oneVideoTask.err = fmt.Errorf("uploadERR: %s", err.Error())
			goto END

		}
		fi, err = fp.Stat()
		if err != nil {
			oneVideoTask.err = fmt.Errorf("uploadERR: %s", err.Error())
			goto END
		}
		para = &map[string]string{"flen": fmt.Sprintf("%d", fi.Size())}
		// para from get to post
		urlOpara = strings.Split(oneVideoTask.Nurl, "?")
		if len(urlOpara) == 2 {
			for _, one = range strings.Split(urlOpara[1], "&") {
				kv = strings.Split(one, "=")
				if len(kv) == 2 {
					(*para)[kv[0]] = kv[1]
				}
			}
		}

		glog.Infoln("uploadmp4: ", oneVideoTask)
		resp, err = upload(oneVideoTask.Nurl, fn+".mp4", para)
		if err != nil {
			oneVideoTask.err = fmt.Errorf("uploadERR: %s", err.Error())
			goto END
		}
		oneVideoTask.nfid = string(resp)

		// 上传新视频缩略图
		jfp, err = os.Open(fn + ".jpg")
		if err != nil {
			glog.Errorln(oneVideoTask, err)
			goto END

		}
		fi, err = jfp.Stat()
		if err != nil {
			glog.Errorln(oneVideoTask, err)
			goto END
		}
		para = &map[string]string{"flen": fmt.Sprintf("%d", fi.Size())}

		glog.Infoln("uploadjpg: ", oneVideoTask)
		jresp, err = upload(oneVideoTask.Nurl, fn+".jpg", para)
		if err != nil {
			glog.Errorln(oneVideoTask, err)
			goto END
		}
		oneVideoTask.nimg = string(jresp)

	END: // 上传完成
		if fp != nil {
			fp.Close()
		}
		if jfp != nil {
			jfp.Close()
		}
		if oneVideoTask.err == nil {
			glog.Infoln("uploadOK: ", oneVideoTask)
		} else {
			glog.Errorf("uploadERR: %v '%s'", oneVideoTask, oneVideoTask.err.Error())
		}
		this.toUpdate <- oneVideoTask
	}
}

func (this *rapperVideo) updateTask() {
	for {
		// 取一个任务
		oneVideoTask := <-this.toUpdate
		
		// 回调
		para := &map[string]string{"type": confJson["tasktype"].(string), "tid": oneVideoTask.Tid, "rid": oneVideoTask.Rid}
		if oneVideoTask.err != nil {
			(*para)["msg"] = base64.StdEncoding.EncodeToString([]byte(oneVideoTask.err.Error()))
		}
		if oneVideoTask.nimg != "" {
			(*para)["img"] = oneVideoTask.nimg
		}
		if oneVideoTask.nfid != "" {
			(*para)["nfid"] = oneVideoTask.nfid
		}

		glog.Infof("callbackTask: %v; para: %v", oneVideoTask, *para)
		_, err := getRequest(oneVideoTask.Callback, para)
		if err != nil {
			olderr := "NULL"
			if oneVideoTask.err != nil {
				olderr = oneVideoTask.err.Error()
			}
			oneVideoTask.err = fmt.Errorf("%s\ncallbackERR: %v", olderr, err)
			glog.Errorln("updateTask callbackERR:", oneVideoTask, oneVideoTask.err.Error())
		}

		// 更新任务状态
		delete(*para, "rid")
		delete(*para, "nfid")
		delete(*para, "img")
		(*para)["stat"] = "1"
		(*para)["name"] = confJson["rappername"].(string)
		if oneVideoTask.err != nil {
			(*para)["stat"] = "-1"
			(*para)["msg"] = base64.StdEncoding.EncodeToString([]byte(oneVideoTask.err.Error()))
		}

		glog.Infoln("updateTask:", oneVideoTask)
		_, err = getRequest(confJson["taskServ"].(string)+"/uptask", para)
		if err != nil {
			glog.Errorln("updateTask updateERR:", oneVideoTask, para, err)
		} else {
			glog.Infoln("updateTaskOK:", oneVideoTask)
		}

		// 删除临时文件
		fn := confJson["tmpdir"].(string) + oneVideoTask.Tid
		os.Remove(fn)
		os.Remove(fn + ".mp4")
		os.Remove(fn + ".jpg")
	}
}
