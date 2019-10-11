package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/liuhengloveyou/easyTask/interfaces"

	"github.com/golang/glog"
)

/*
{
 "fid":"900499",
 "flen":15070734,
 "type":"wmv",
 "url":"http://www.upload.cn/client14081804173567259325.wmv",
 "nurl":"http://192.168.1.2:8081/ffaceup",
 "callback":"http://www.space.cn/opus/updateNfid.html",
 "oldpath":"domain.cn:/file/name"
}
*/
type videoTaskInfo1 struct {
	Flen     int32  // 文件长度
	Fid      string // 文件ID
	Type     string // 文件类型
	Url      string // 文件下载URL
	Nurl     string // 文件长度
	Callback string // 任务回调
	Oldpath  string // 老路径
	Tid      string // 任务ID
	Rid      string // 记录ID
	nfid     string // 转码处理后的文件ID
	nimg     string // 缩略图
	nimg1    string // 缩略图1
	err      error  // 处理过程中的错误信息
}

type rapperVideo1 struct {
	no          int64
	toTranscode chan *videoTaskInfo1
	toUpload    chan *videoTaskInfo1
	toUpdate    chan *videoTaskInfo1
}

func init() {
	interfaces.Register("video1", NewRapperVideo1)
}

func NewRapperVideo1() rapper {
	return &rapperVideo1{
		no:          time.Now().UnixNano(),
		toTranscode: make(chan *videoTaskInfo1, 5),
		toUpload:    make(chan *videoTaskInfo1, 5),
		toUpdate:    make(chan *videoTaskInfo1, 5)}
}

func (this *rapperVideo1) run() {
	uploadUrlConf := ""
	callbackUrlConf := ""
	if val, ok := confJson["uploadUrl"]; ok == true && val.(string) != "" {
		uploadUrlConf = val.(string)
	}
	if val, ok := confJson["callbackUrl"]; ok == true && val.(string) != "" {
		callbackUrlConf = val.(string)
	}

	go this.transcode()
	go this.uploadMp4()
	go this.updateTask()

	for {
		// 取一个任务
		oneTask := oneTask()
		oneVideoTask := &videoTaskInfo1{Tid: oneTask.Tid, Rid: oneTask.Rid}
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

		// 覆盖配置参数
		if uploadUrlConf != "" {
			oneVideoTask.Nurl = uploadUrlConf
		}
		if callbackUrlConf != "" {
			oneVideoTask.Callback = callbackUrlConf
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

func (this *rapperVideo1) transcode() {
	const TRANFMT = "export LANGUAGE=en_US;ffmpeg -y -v error -i %s -movflags +faststart -vcodec libx264 -profile:v main -level 3 -b:v 512k -acodec libvo_aacenc -b:a 128k -vf 'movie=watermark.png[logo];[in][logo]overlay=main_w-overlay_w-5:10[out]' %s.mp4"
	const TRANFMT1 = "export LANGUAGE=en_US;ffmpeg -y -v error -i %s -movflags +faststart -vcodec libx264 -profile:v main -level 3 -b:v 800k -acodec libvo_aacenc -b:a 128k -s 800x450 -vf 'movie=watermark.png[logo];[in][logo]overlay=main_w-overlay_w-5:10[out]' %s.mp4"

	const THUMFMT = "export LANGUAGE=en_US;ffmpeg -y -v error -i %s.mp4 -y -f image2 -ss 8.0  -vframes 1 -s 120x80 %s.jpg"
	const THUMFM1 = "export LANGUAGE=en_US;ffmpeg -y -v error -i %s.mp4 -y -f image2 -ss 8.0  -vframes 1 -s 320x240 %s_1.jpg"
	const THUMFM2 = "export LANGUAGE=en_US;ffmpeg -y -v error -i %s_1.jpg -y -f image2 -vf crop=320:170:0:0 -s 200x112 %s_1.jpg"

	for {
		// 取一个任务
		oneVideoTask := <-this.toTranscode
		fn := confJson["tmpdir"].(string) + oneVideoTask.Tid

		transcodeCmd := fmt.Sprintf(TRANFMT, fn, fn)
		thumbnailCmd := fmt.Sprintf(THUMFMT, fn, fn)
		thumbnailCm1 := fmt.Sprintf(THUMFM1, fn, fn)
		thumbnailCm2 := fmt.Sprintf(THUMFM2, fn, fn)

		// 取原视频大小
		cmd := fmt.Sprintf(`ffmpeg -i %s 2>&1 | grep 'Stream.*Video.*' | sed -r 's/^.*, ([0-9]{3,4}x[0-9]{3,4}).*/\1/g'`, fn)
		fsize, _ := exec.Command("/bin/bash", "-c", cmd).Output()
		if len(fsize) > 3 {
			wh := strings.Split(strings.Replace(string(fsize), "\n", "x", -1), "x")
			if len(wh) >= 2 {
				wi, errw := strconv.Atoi(wh[0])
				hi, errh := strconv.Atoi(wh[1])
				if errw == nil && errh == nil && (wi > 800 || hi > 600) {
					transcodeCmd = fmt.Sprintf(TRANFMT1, fn, fn)
					glog.Infof("%s transcodeCMD: %d. cmd='%s';", oneVideoTask.toString(), this.no, transcodeCmd)
				}
			}
		}

		// 转码
		glog.Infoln("transcode: ", this.no, oneVideoTask.toString(), transcodeCmd)
		_, err := exec.Command("/bin/bash", "-c", transcodeCmd).Output()
		if err != nil {
			oneVideoTask.err = fmt.Errorf("transcodeERR: cmd='%s'; err='%s'", transcodeCmd, err.Error())
			goto END
		}

		// 取缩略图
		glog.Infoln("thumbnail: ", this.no, oneVideoTask.toString(), thumbnailCmd)
		_, err = exec.Command("/bin/bash", "-c", thumbnailCmd).Output()
		if err != nil {
			glog.Errorln(err, thumbnailCmd)
		}

		// 取缩略图1
		glog.Infoln("thumbnail: ", this.no, oneVideoTask.toString(), thumbnailCm1)
		_, err = exec.Command("/bin/bash", "-c", thumbnailCm1).Output()
		if err != nil {
			glog.Errorln(err, thumbnailCm1)
		}

		// 取缩略图2
		glog.Infoln("thumbnail: ", this.no, oneVideoTask.toString(), thumbnailCm2)
		_, err = exec.Command("/bin/bash", "-c", thumbnailCm2).Output()
		if err != nil {
			glog.Errorln(err, thumbnailCm2)
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

func (this *rapperVideo1) uploadMp4() {
	for {
		var resp, jresp []byte
		para := map[string]string{}

		// 取一个任务
		oneVideoTask := <-this.toUpload

		// 上传新视频
		fi, err := os.Stat(confJson["tmpdir"].(string) + oneVideoTask.Tid + ".mp4")
		if err != nil {
			oneVideoTask.err = fmt.Errorf("uploadERR: %s", err.Error())
			goto END
		}

		para = map[string]string{"flen": fmt.Sprintf("%d", fi.Size()), "oldpath": fmt.Sprintf("%s.mp4", oneVideoTask.Oldpath)}
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
		para = map[string]string{"flen": fmt.Sprintf("%d", fi.Size()), "oldpath": fmt.Sprintf("%s.jpg", oneVideoTask.Oldpath)}
		glog.Infoln("uploadjpg: ", this.no, oneVideoTask.toString(), para)
		jresp, err = upload(oneVideoTask.Nurl, confJson["tmpdir"].(string)+oneVideoTask.Tid+".jpg", &para)
		if err != nil {
			glog.Errorln(oneVideoTask.toString(), err.Error())
			goto END
		}
		oneVideoTask.nimg = strings.Trim(string(jresp), "\n ")

		// 上传新视频缩略图1
		fi, err = os.Stat(confJson["tmpdir"].(string) + oneVideoTask.Tid + "_1.jpg")
		if err != nil {
			glog.Errorln(oneVideoTask.toString(), err)
			goto END
		}
		para = map[string]string{"flen": fmt.Sprintf("%d", fi.Size()), "oldpath": fmt.Sprintf("%s_1.jpg", oneVideoTask.Oldpath)}
		glog.Infoln("uploadjpg: ", this.no, oneVideoTask.toString(), para)
		jresp, err = upload(oneVideoTask.Nurl, confJson["tmpdir"].(string)+oneVideoTask.Tid+"_1.jpg", &para)
		if err != nil {
			glog.Errorln(oneVideoTask.toString(), err.Error())
			goto END
		}
		oneVideoTask.nimg1 = strings.Trim(string(jresp), "\n ")

	END: // 上传完成
		if oneVideoTask.err == nil {
			glog.Warningln("uploadOK: ", this.no, oneVideoTask.toString())
		} else {
			glog.Errorln("uploadERR: ", this.no, oneVideoTask.toString())
		}
		this.toUpdate <- oneVideoTask
	}
}

func (this *rapperVideo1) updateTask() {
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
		if oneVideoTask.nimg1 != "" {
			para["img1"] = oneVideoTask.nimg1
		}

		olderr := "NULL"
		if oneVideoTask.err != nil {
			olderr = oneVideoTask.err.Error()
		}
		glog.Warningln("callbackTask: ", this.no, oneVideoTask.toString(), para)
		body, err := getRequest(oneVideoTask.Callback, &para)
		if err != nil {
			oneVideoTask.err = fmt.Errorf("%s\ncallbackERR: %s", olderr, err.Error())
			glog.Errorln("updateTask callbackERR:", oneVideoTask.toString())
		} else if string(body) != "true" {
			oneVideoTask.err = fmt.Errorf("%s\ncallbackERR: %s", olderr, body)
			glog.Errorln("updateTask callbackERR:", oneVideoTask.toString(), body)
		}

		// 更新任务状态
		delete(para, "rid")
		delete(para, "nfid")
		delete(para, "img")
		delete(para, "img1")
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
		os.Remove(fn + "_1.jpg")
	}
}

func (this *videoTaskInfo1) toString() string {
	errstr := ""
	if this.err != nil {
		errstr = this.err.Error()
	}
	return fmt.Sprintf("{flen: %d; fid: %s; type: %s; url: %s; nurl: %s; callback: %s; tid: %s; rid: %s; nfid: %s; nimg: %s; nimg1: %s; err: %s}",
		this.Flen, this.Fid, this.Type, this.Url, this.Nurl, this.Callback, this.Tid, this.Rid, this.nfid, this.nimg, this.nimg1, errstr)
}
