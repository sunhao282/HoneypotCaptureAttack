package docker

import (
	"fmt"
	"honey/src/util/log"
	"io/ioutil"
	"net/http"
	"strings"
	"unsafe"
)

func Start(addr string) {
	dockermux := http.NewServeMux()

	dockermux.HandleFunc("/", FakedockerBanner)
	dockermux.HandleFunc("/_ping", Fakedocker)
	dockermux.HandleFunc("/version", Fakeversiondocker)
	dockermux.HandleFunc("/info", Fakeinfodocker)
	dockermux.HandleFunc("/v1.6/version", Fakeversiondocker)
	dockermux.HandleFunc("/v1.24/containers/json", Fakecontainersdocker)
	// Start the server   监听端口
	http.ListenAndServe(addr, dockermux)
}
func Fakecontainersdocker(w http.ResponseWriter, r *http.Request) {
	printInfo(r)
	response := `[]`
	WritejsonResponse(w, response)
	return
}
func Fakeinfodocker(w http.ResponseWriter, r *http.Request) {
	printInfo(r)
	response := `{"ID":"W7DU:U7FW:AFFX:R744:H6RI:NVUH:LJKM:KV6F:P72S:HL3T:SGET:RYZI","Containers":1,"ContainersRunning":0,"ContainersPaused":0,"ContainersStopped":1,"Images":2,"Driver":"overlay2","DriverStatus":[["Backing Filesystem","xfs"],["Supports d_type","true"],["Native Overlay Diff","true"]],"SystemStatus":null,"Plugins":{"Volume":["local"],"Network":["bridge","host","ipvlan","macvlan","null","overlay"],"Authorization":null,"Log":["awslogs","fluentd","gcplogs","gelf","journald","json-file","local","logentries","splunk","syslog"]},"MemoryLimit":true,"SwapLimit":true,"KernelMemory":true,"KernelMemoryTCP":true,"CpuCfsPeriod":true,"CpuCfsQuota":true,"CPUShares":true,"CPUSet":true,"PidsLimit":true,"IPv4Forwarding":false,"BridgeNfIptables":false,"BridgeNfIp6tables":false,"Debug":false,"NFd":26,"OomKillDisable":true,"NGoroutines":41,"SystemTime":"2020-03-03T18:28:06.818532528-08:00","LoggingDriver":"json-file","CgroupDriver":"cgroupfs","NEventsListener":0,"KernelVersion":"3.10.0-957.el7.x86_64","OperatingSystem":"CentOS Linux 7 (Core)","OSType":"linux","Architecture":"x86_64","IndexServerAddress":"https://index.docker.io/v1/","RegistryConfig":{"AllowNondistributableArtifactsCIDRs":[],"AllowNondistributableArtifactsHostnames":[],"InsecureRegistryCIDRs":["127.0.0.0/8"],"IndexConfigs":{"docker.io":{"Name":"docker.io","Mirrors":[],"Secure":true,"Official":true}},"Mirrors":[]},"NCPU":1,"MemTotal":1019797504,"GenericResources":null,"DockerRootDir":"/var/lib/docker","HttpProxy":"","HttpsProxy":"","NoProxy":"","Name":"localhost.localdomain","Labels":[],"ExperimentalBuild":false,"ServerVersion":"19.03.5","ClusterStore":"","ClusterAdvertise":"","Runtimes":{"runc":{"path":"runc"}},"DefaultRuntime":"runc","Swarm":{"NodeID":"","NodeAddr":"","LocalNodeState":"inactive","ControlAvailable":false,"Error":"","RemoteManagers":null},"LiveRestoreEnabled":false,"Isolation":"","InitBinary":"docker-init","ContainerdCommit":{"ID":"b34a5c8af56e510852c35414db4c1f4fa6172339","Expected":"b34a5c8af56e510852c35414db4c1f4fa6172339"},"RuncCommit":{"ID":"3e425f80a8c931f88e6d94a8c831b9d5aa481657","Expected":"3e425f80a8c931f88e6d94a8c831b9d5aa481657"},"InitCommit":{"ID":"fec3683","Expected":"fec3683"},"SecurityOptions":["name=seccomp,profile=default"],"Warnings":["WARNING: API is accessible on http://0.0.0.0:2375 without encryption.\n         Access to the remote API is equivalent to root access on the host. Refer\n         to the 'Docker daemon attack surface' section in the documentation for\n         more information: https://docs.docker.com/engine/security/security/#docker-daemon-attack-surface","WARNING: IPv4 forwarding is disabled","WARNING: bridge-nf-call-iptables is disabled","WARNING: bridge-nf-call-ip6tables is disabled"]}`
	WritejsonResponse(w, response)
	return
}
func Fakeversiondocker(w http.ResponseWriter, r *http.Request) {
	printInfo(r)
	response := `{"Platform":{"Name":"Docker Engine - Community"},"Components":[{"Name":"Engine","Version":"19.03.5","Details":{"ApiVersion":"1.40","Arch":"amd64","BuildTime":"2019-11-13T07:24:18.000000000+00:00","Experimental":"false","GitCommit":"633a0ea","GoVersion":"go1.12.12","KernelVersion":"3.10.0-957.el7.x86_64","MinAPIVersion":"1.12","Os":"linux"}},{"Name":"containerd","Version":"1.2.10","Details":{"GitCommit":"b34a5c8af56e510852c35414db4c1f4fa6172339"}},{"Name":"runc","Version":"1.0.0-rc8+dev","Details":{"GitCommit":"3e425f80a8c931f88e6d94a8c831b9d5aa481657"}},{"Name":"docker-init","Version":"0.18.0","Details":{"GitCommit":"fec3683"}}],"Version":"19.03.5","ApiVersion":"1.40","MinAPIVersion":"1.12","GitCommit":"633a0ea","GoVersion":"go1.12.12","Os":"linux","Arch":"amd64","KernelVersion":"3.10.0-957.el7.x86_64","BuildTime":"2019-11-13T07:24:18.000000000+00:00"}`
	WritejsonResponse(w, response)
	return
}
func Fakedocker(w http.ResponseWriter, r *http.Request) {
	printInfo(r)
	response := `{"_ping":"ok"}`
	WriteResponse(w, response)
	return
}
func FakedockerBanner(w http.ResponseWriter, r *http.Request) {
	printInfo(r)
	response := `{"message":"page not found"}`
	WritejsonResponse(w, response)
	return
}

func WritejsonResponse(w http.ResponseWriter, d string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Server", "Docker/19.03.5 (linux)")
	w.Header().Set("Ostype", "linux")
	w.Header().Set("Api-Version", "1.40")
	w.WriteHeader(200)
	w.Write([]byte(d))
	return
}
func printInfo(r *http.Request) {
	info := "URL:" + r.URL.String() + "&&Method:" + r.Method + "&&RemoteAddr:" + r.RemoteAddr
	arr := strings.Split(r.RemoteAddr, ":")
	if strings.EqualFold(r.Method, "post") {
		respBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err.Error())
		}
		strr := (*string)(unsafe.Pointer(&respBytes))
		log.Pr("docker", arr[0], "Url", info, strip(*strr, "/ \n\r\t"))
	} else {
		//arr := strings.Split(r.RemoteAddr, ":")
		log.Pr("docker", arr[0], "Url", info)
	}
}
func strip(s_ string, chars_ string) string {
	s, chars := []rune(s_), []rune(chars_)
	length := len(s)
	max := len(s) - 1
	l, r := true, true //标记当左端或者右端找到正常字符后就停止继续寻找
	start, end := 0, max
	tmpEnd := 0
	charset := make(map[rune]bool) //创建字符集，也就是唯一的字符，方便后面判断是否存在
	for i := 0; i < len(chars); i++ {
		charset[chars[i]] = true
	}
	for i := 0; i < length; i++ {
		if _, exist := charset[s[i]]; l && !exist {
			start = i
			l = false
		}
		tmpEnd = max - i
		if _, exist := charset[s[tmpEnd]]; r && !exist {
			end = tmpEnd
			r = false
		}
		if !l && !r {
			break
		}
	}
	if l && r { // 如果左端和右端都没找到正常字符，那么表示该字符串没有正常字符
		return ""
	}
	return string(s[start : end+1])
}
func WriteResponse(w http.ResponseWriter, d string) {
	w.Header().Set("Api-Version", "1.40")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	//w.Header().Set("Content-Length", "2")
	w.Header().Set("Content-Type", "test/plain; charset=utf-8")
	w.Header().Set("Docker-Experimental", "false")
	w.Header().Set("Ostype", "linux")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Server", "Docker/19.03.5 (linux)")
	w.WriteHeader(200)
	//redate:=gojson.Json(d)
	//byteDate,err:=json.Marshal(redate.Getdata())
	//if err != nil {
	//	fmt.Println("json 解析错误" )
	w.Write([]byte(d))
	return
}
