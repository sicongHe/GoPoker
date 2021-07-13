package server

import "fmt"

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{},[]Player{}}
}

type InMemoryPlayerStore struct {
	store map[string]int
	league []Player
}

func (imp InMemoryPlayerStore) GetPlayerScore(name string) string{
	return fmt.Sprintf("%d",imp.store[name])
}
func (imp *InMemoryPlayerStore) RecordWin(name string) {
	imp.store[name]++
}

func (imp *InMemoryPlayerStore) GetLeague() []Player {
	var ret []Player
	for name,wins := range imp.store {
		ret = append(ret, Player{name, wins})
	}
	return ret
}

