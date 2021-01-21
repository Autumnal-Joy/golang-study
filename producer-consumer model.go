package main

import "fmt"

type orderInfo struct {
	id int
}

func producer(ch chan<- orderInfo) {
	for i := 0; i < 10; i++ {
		order := orderInfo{id: i}
		ch <- order
	}
	close(ch)
}

func consumer(ch <-chan orderInfo) {
	for v := range ch {
		fmt.Println("订单id为: ", v.id)
	}
}

func main() {
	ch := make(chan orderInfo, 5)
	go producer(ch)
	consumer(ch)
}
