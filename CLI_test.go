package poker_test

import (
	poker "github.com/siconghe/MyServer"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Run("nancy wins", func(t *testing.T) {
		in := strings.NewReader("nancy wins!\n")
		playerStore := &poker.StubPlayerStore{}
		cli := poker.NewCLI(playerStore,in)
		cli.PlayPoker()
		assertPlayerWins(t,playerStore,"nancy")
	})
	t.Run("tom wins", func(t *testing.T) {
		in := strings.NewReader("tom wins!\n")
		playerStore := &poker.StubPlayerStore{}
		cli := poker.NewCLI(playerStore,in)
		cli.PlayPoker()
		assertPlayerWins(t,playerStore,"tom")
	})
}

func assertPlayerWins(t *testing.T,playerStore *poker.StubPlayerStore,name string) {
	if len(playerStore.WinCalls) !=1 {
		t.Fatal("expected a win call but didn't get any")
	}
	got := playerStore.WinCalls[0]
	want := name
	poker.AssertString(t,got,want)
}
