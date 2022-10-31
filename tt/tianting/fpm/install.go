package fpm

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"tianting/comm"
)

var log comm.ConsoleLogger

func Install(softwareName, mode, version, path, config string) (string, string) {
	_, err := os.Stat(path)
	if err != nil {
		log.Fatal("目录不存在%v", err)
	}
	if os.IsNotExist(err) {
		if err = os.Mkdir(path, 0755); err != nil {
			log.Debug("%v 目录创建失败！目录信息:%v", path)
		} else {
			log.Debug("%v目录不存在，已创建:", path)
		}
	}

	file_all_path, b := Wget(softwareName, version)
	log.Debug("wget func respose is:%v,%v", file_all_path, b)
	switch b {
	case true:
		rpm_path := Makerpms(file_all_path, path, softwareName, version, true)
		return file_all_path, rpm_path
	case false:
		if file_all_path == "ERROR" {
			log.Fatal("文件下载失败,judge:%v", file_all_path)
			return "download error", "下载失败"
		} else {
			rpm_path := Makerpms(file_all_path, path, softwareName, version, false)
			return file_all_path, rpm_path
		}
	}
	return file_all_path, "rpm_path"
}

//fpm -s dir -t rpm -n nginx -v 1.6.3 -d 'pcre-devel,openssl-devel' --post-install /server/scripts/nginx_rpm.sh -f /application/nginx-1.6.3/
func Makerpms(file_all_path, soft_dir, softwareName, softwareVersion string, b bool) string {

	cmd0 := exec.Command("/bin/bash", "-c", "/usr/bin/tar"+" xf "+file_all_path+" -C "+soft_dir)
	log.Debug("压缩包解压中...")
	out0, err := cmd0.CombinedOutput()
	if err != nil {
		log.Error("查看%v没有解压成功，错误是%v", file_all_path, string(out0))
	}
	log.Debug("解压完成，rpm包制作中")
	cmd := exec.Command("/bin/bash", "-c", "/usr/local/rvm/gems/ruby-2.7.2/bin/fpm"+" -s"+" dir"+" -t"+" rpm"+" -p /data/workspace/tt/src/tianting/fpm/output "+"--post-install instance/"+softwareName+"_post_install.sh"+" -n "+softwareName+" -v "+softwareVersion+" -f "+soft_dir+" instance/"+softwareName+".service")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	//{:timestamp=>"2022-05-28T11:46:45.014426+0800", :message=>"Created package", :path=>"/data/workspace/tt/src/tianting/fpm/output/redis-6.2.2-1.x86_64.rpm"}
	//截取切片的结果全路径
	str := string(out)
	ret_str := strings.Split(str, "\"")
	//log.Debug("rpm制作的 --------完整信息:%v", str)
	//资源是否第一次制作，rpm的返回值不同，导致取rpm的路径时需要做判断
	switch b {
	case true:
		log.Debug("使用本地包制作完成.rpm存储路径:%v", ret_str[len(ret_str)-2])
		//log.Debug("制作完成.rpm存储路径: --------完整信息:%v", str)
		return ret_str[len(ret_str)-2]
	case false:
		log.Debug("使用下载包制作完成.rpm存储路径: %v", ret_str[len(ret_str)-2])
		//log.Debug("制作完成.rpm存储路径: --------完整信息:%v", str)
		return ret_str[len(ret_str)-2]
	}

	return ret_str[len(ret_str)-2]
}
