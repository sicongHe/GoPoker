package main

import (
	"github.com/siconghe/MyServer"
	"log"
	"net/http"
)

const dbFileName = "./cmd/webserver/game.db.json"
func main() {
	store, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatalf("服务器初始化失败 %s", err.Error())
	}
	server := poker.NewPlayerServer(store)
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}