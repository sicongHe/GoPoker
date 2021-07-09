package server

import (
	"fmt"
	"net/http"
)

type PlayerServer struct {
	Store PlayerStore
}


type PlayerStore interface {
	GetPlayerScore(name string) string
	RecordWin(name string)
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

type InMemoryPlayerStore struct {
	store map[string]int
}

func (imp InMemoryPlayerStore) GetPlayerScore(name string) string{
	return fmt.Sprintf("%d",imp.store[name])
}
func (imp *InMemoryPlayerStore) RecordWin(name string) {
	imp.store[name]++
}
func (ps *PlayerServer)ServeHTTP(rw http.ResponseWriter, rq *http.Request){
	switch rq.Method {
	case http.MethodGet:
		ps.showCode(rw,rq)
	case http.MethodPost:
		ps.processWin(rw,rq)
	}
}

func (ps *PlayerServer)processWin(rw http.ResponseWriter,rq *http.Request) {
	player := rq.URL.Path[len("/players/"):]
	ps.Store.RecordWin(player)
	rw.WriteHeader(http.StatusAccepted)
}

func (ps *PlayerServer)showCode(rw http.ResponseWriter,rq *http.Request) {
	player := rq.URL.Path[len("/players/"):]
	score := ps.Store.GetPlayerScore(player)
	if score == "" {
		rw.WriteHeader(404)
	}
	fmt.Fprintf(rw,score)
}




