package setting

import (
	"fmt"
	"honey/src/core/exec"
	"honey/src/core/exec/protocol/docker"
	"honey/src/util/conf"
)

func Help() {
	exec.Execute("clear") //[/bin/bash -c clear]

	fmt.Println("WS Active Attack Honeypot Fishing")
	fmt.Println("------------------------------------------------")
	fmt.Println("")
	fmt.Println("   -r,   Start up service")
	fmt.Println("   -v,   Version")
	fmt.Println("   -h,   Help")
	fmt.Println("")
	fmt.Println("------------------------------------------------")
	fmt.Println("")
}

//func Init() {
//	fmt.Println("init...")
//}
func Version() {
	fmt.Println("1.0")
}
func Run() {
	// 启动 docker 蜜罐
	//fmt.Println("启动docker蜜罐中.....")
	dockerStatus := conf.Get("docker", "status")
	//fmt.Println("dockerStatus: ", dockerStatus)

	// 判断 docker 蜜罐 是否开启
	if dockerStatus == "1" {
		dockerAddr := conf.Get("docker", "addr")
		//fmt.Println("dockerAddr:", dockerAddr)
		go docker.Start(dockerAddr) //go关键字 开启并行
	}
	//fmt.Println("启动docker蜜罐成功！")
}
