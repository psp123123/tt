package main

import (
	"fmt"

	"runtime"

)

func main() {
	sysType:=runtime.GOOS
	if sysType=="linux" {
		fmt.Printf("当前系统是：%v\n",sysType)
	} else {
		fmt.Printf("当前系统是：%v\n",sysType)
	}
}