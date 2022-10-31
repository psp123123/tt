package fpm

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"tianting/conn"
	"tianting/tools"
)

func Wget(name, Version string) (string, bool) {
	base_dir := tools.JudgeOS()
	work_download_dir := base_dir + "src/tianting/fpm/resources"
	downLoadName := name + "-" + Version + ".tar.gz"
	_, err := os.Stat(work_download_dir + "/" + downLoadName)
	if err == nil {
		log.Debug("文件已存在将使用本地文件,路径是：%v", work_download_dir+"/"+downLoadName)
		return work_download_dir + "/" + downLoadName, true
	}
	if os.IsNotExist(err) {
		//return downLoadName +"is not exist"
		log.Debug("文件不存在，将网上下载")
		Url := conn.GetDownLoadURL(name)
		log.Debug("%v开始下载...", Url+downLoadName)
		res, err := http.Get(Url + downLoadName)
		if err != nil {
			log.Error("url is valid:%v", err)
			return "error", false
		}
		log.Debug("Internet resource weither exist: %v", res.Status)
		if res.Status != "200 OK" {
			log.Fatal("所需资源不存在,Internet resource is not exist, Please check URL! ")
			return "ERROR", false
		}

		defer res.Body.Close()
		downLoadFileNmae, err := os.Create(work_download_dir + "/" + downLoadName)
		if err != nil {
			log.Error("主机创建文件失败:%v", err)
		}
		writer := bufio.NewWriter(downLoadFileNmae)
		written, _ := io.Copy(writer, res.Body)
		log.Debug("total lenth :%v,下载完成.", written)
		return work_download_dir + "/" + downLoadName, false
	}
	return work_download_dir + "/" + downLoadName, false

}
