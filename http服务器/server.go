package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// 注册回调函数
	http.HandleFunc("/", handler)

	// 监听端口
	err := http.ListenAndServe("127.0.0.1:8888", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(request.Method, request.URL)

	// 解析资源地址
	src, _ := filepath.Abs("static" + request.URL.String())

	// 获取资源类型
	stat, err := os.Stat(src)
	if err != nil {
		fmt.Println(err)
		writer.Write([]byte("resource not found"))
		return
	}

	if stat.IsDir() {
		fileList, _ := ioutil.ReadDir(src)

		for _, f := range fileList {
			info := fmt.Sprintf("%s %dbytes\n", f.Name(), f.Size())
			writer.Write([]byte(info))
		}
	} else {
		file, _ := os.Open(src)
		defer file.Close()

		buf := make([]byte, 4096)
		for {
			n, _ := file.Read(buf)
			if n == 0 {
				break
			}
			writer.Write(buf[:n])
		}
	}

}
