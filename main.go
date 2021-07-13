package main

import (
	"fmt"
	server2 "github.com/siconghe/MyServer/server"
	"log"
	"net/http"
)


func main() {
	fmt.Println("开始开发简易HTTP服务器！")
	server := server2.NewPlayerServer(server2.NewInMemoryPlayerStore())
	if err := http.ListenAndServe(":5000",server); err!= nil {
		log.Fatal("服务器启动失败")
	}
}