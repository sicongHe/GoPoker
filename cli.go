package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	Store PlayerStore
	In *bufio.Scanner
}

func NewCLI(store PlayerStore,in io.Reader) *CLI {
	return &CLI{
		store,
		bufio.NewScanner(in),
	}
}

func (cli *CLI) readline() string {
	cli.In.Scan()
	return cli.In.Text()
}

func extractWinner(input string) string {
	return strings.Replace(input," wins!","",-1)
}

func (cli *CLI) PlayPoker() {
	userInput := cli.readline()

	cli.Store.RecordWin(extractWinner(userInput))
}
