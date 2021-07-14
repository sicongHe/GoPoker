package main

import (
	"fmt"
	"github.com/siconghe/MyServer"
	"log"
	"os"
)

const dbFileName = "./cmd/cli/game.db.json"
func main() {
	fmt.Println("XXX wins!")
	store, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatalf("服务器初始化失败 %s", err.Error())
	}
	game := poker.NewCLI(store,os.Stdin)
	game.PlayPoker()
}

