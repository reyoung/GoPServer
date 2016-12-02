package main

import "github.com/reyoung/GoPServer/service"

func main() {
	service.New("tcp://127.0.0.1:9000")
}
