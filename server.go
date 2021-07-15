// Package poker
//为WebServer和Cli应用程序提供支持的库程序
//
///*
package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Player struct {
	Name string
	Wins int
}

type PlayerServer struct {
	Store PlayerStore
	http.Handler
}


type PlayerStore interface {
	GetPlayerScore(name string) string
	RecordWin(name string)
	GetLeague() League
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	ps := new(PlayerServer)
	ps.Store = store
	router := http.NewServeMux()
	router.Handle("/league",http.HandlerFunc(ps.leagueHandler))
	router.Handle("/players/",http.HandlerFunc(ps.playersHandler))
	ps.Handler = router
	return ps
}



func (ps *PlayerServer) leagueHandler(rw http.ResponseWriter, rq *http.Request){
	leagueTable := ps.Store.GetLeague()
	json.NewEncoder(rw).Encode(leagueTable)
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("content-type","application-json")
}



func (ps *PlayerServer) GetLeagueTable() []Player {
	return ps.Store.GetLeague()
}



func (ps *PlayerServer) playersHandler(rw http.ResponseWriter, rq *http.Request){
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




