package tools

import (
	"runtime"
	"tianting/comm"
)

var log comm.ConsoleLogger

func JudgeOS() string {
	//base_dir :=""
	sysType := runtime.GOOS
	if sysType == "linux" {

		base_dir := "/data/workspace/tt/"
		log.Debug("当前系统是：%v,路径是：%v\n", sysType, base_dir)
		return base_dir
	} else {
		base_dir := ""
		log.Debug("当前系统是：%v,路径是：%v\n", sysType, base_dir)
		return base_dir

	}
	//return base_dir
}
