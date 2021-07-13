package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type League []Player

type FileSystemStore struct {
	Database io.ReadWriteSeeker
	league League
}

func NewFileSystemStore(database io.ReadWriteSeeker) *FileSystemStore {
	database.Seek(0, 0)
	league,_ := NewLeague(database)
	return &FileSystemStore{database,league}
}

func (l League)Find(name string) *Player {
	for i,player := range l {
		if player.Name == name {
			return &l[i]
		}
	}
	return nil
}

func (f *FileSystemStore) GetPlayerScore(name string) string {
	var ret string
	for _,player := range f.league {
		if player.Name == name {
			ret = fmt.Sprintf("%d",player.Wins)
		}
	}
	return ret
}

func (f *FileSystemStore) RecordWin(name string) {
	player := f.league.Find(name)
	if player == nil {
		f.league = append(f.league,Player{name,1})
	} else{
		player.Wins++
	}
	f.Database.Seek(0,0)
	json.NewEncoder(f.Database).Encode(&(f.league))
}

func (f *FileSystemStore) GetLeague() League {
	return f.league
}

func NewLeague(reader io.Reader)(League, error) {
	var league League
	err := json.NewDecoder(reader).Decode(&league)
	if err!= nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}
	return league,err
}

func (f *FileSystemStore) testFunc() {
	context.TODO()
}

