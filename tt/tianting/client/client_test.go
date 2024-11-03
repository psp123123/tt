package client

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"time"
)

func client_ping() interface{} {

	counter := 0
	for {

		conn, err := net.DialTimeout("tcp", "47.93.37.151:80", 3000)
		if err != nil {
			counter++
			fmt.Printf("网络及端口不通：%v\n,第%v次尝试", err, counter)
			time.Sleep(time.Second)
			//return err
		} else {
			fmt.Printf("conn is %v\n", conn)
			time.Sleep(1000 * time.Millisecond)
			return conn

		}
	}

}

func client_get() {
	counter := 0
	for {

		var resp, err = http.Get("http://47.93.37.151:80/svc_names")
		if err != nil {
			fmt.Printf("访问服务名称列表失败：因为：%v\n", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("解析数据失败：%v\n", err)
		}
		fmt.Printf("得到数据：%s\n", body)
		counter++
		time.Sleep(1000 * time.Millisecond)
	}
}

func TestPostConfig(softwareName, mode, way, version, path, buildArgs, config string) []string {
	//var list struct{
	//	softwareName string
	//	mode	string
	//	way	string
	//	version string
	//	path string
	//	buildArgs string
	//	config string
	//}
	//var m1 []string
	//m1 := make([]string)
	//m1[list.softwareName] = softwareName
	//m1[list.mode] = mode
	//m1[list.way] = way
	//m1[list.version] = version
	//m1[list.path] = path
	//m1[list.buildArgs] = buildArgs
	//m1[list.config] = config
	//m1[0] = softwareName
	//m1[1] = mode
	//m1[0] = way
	//m1[0] = version
	//m1[4] = path
	//m1[5] = buildArgs
	//m1[6] = config
	m1 := []string{softwareName, mode, way, version, path, buildArgs, config}

	return m1
}

func client_post() {
	counter := 0
	//Data := TestPostConfig("a","b","c","d","e","f","g")
	//var SData []string
	//
	//for _,v := range Data {
	//	SData[counter] = v
	//	counter ++
	//}
	for {

		//var resp, err = http.PostForm("http://127.0.0.1:3000/submit", url.Values{"softwareName":{"a"},"mode" :{"b"},"way":{"c"} , "version":{"d"} ,"path":{"e1111111111"} ,"buildArgs":{"f"} ,"config":{"g"}})
		var resp, err = http.PostForm("http://47.93.37.151/submit", url.Values{"softwareName": {"a"}, "mode": {"111"}, "way": {"c"}, "version": {"d"}, "path": {"e1111111111"}, "buildArgs": {"f"}, "config": {"g"}})
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		counter++
		if err != nil {
			fmt.Printf("httpPost receive context is err:%v", err)
		}
		fmt.Println(string(body), counter)

		time.Sleep(1000 * time.Millisecond)
	}
}
func client() {

	client_ping()
	client_post()

}

//func install(get_url string,server string,ver string)  {
//	os.
//
//}
