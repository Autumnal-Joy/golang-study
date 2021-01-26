package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	urls := []string{"https://www.baidu.com/", "https://www.bilibili.com/"}
	quit := make(chan bool)
	for i, url := range urls {
		go func(i int, url string) { // 获取服务器应答包
			resp, err := http.Get(url)
			if err != nil {
				log.Fatalln("http.Get err: ", err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln("ioutil.ReadAll err: ", err)
			}
			err = ioutil.WriteFile(fmt.Sprintf("%d.html", i), body, 0666)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(url + " finished")
			quit <- true
		}(i, url)
	}
	for i := 0; i < len(urls); i++ {
		<-quit
	}
}
