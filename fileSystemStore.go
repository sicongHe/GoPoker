package poker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

type League []Player

type FileSystemStore struct {
	Database *json.Encoder
	league   League
}

func FileSystemPlayerStoreFromFile(name string) (*FileSystemStore, error) {
	db, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening %s %v", name, err)
	}
	store ,err:= NewFileSystemStore(db)
	return store,err
}

func NewFileSystemStore(database *os.File) (*FileSystemStore,error) {
	database.Seek(0, 0)
	info,err := database.Stat()
	if err!= nil {
		return nil,fmt.Errorf("初始化FileSystem失败，错误信息: %s",err)
	}
	if info.Size() == 0 {
		database.Write([]byte("[]"))
		database.Seek(0, 0)
	}
	league,err := NewLeague(database)
	if err!= nil {
		return nil,fmt.Errorf("初始化FileSystem失败，错误信息: %s",err)
	}
	return &FileSystemStore{json.NewEncoder(&Tape{database}),league},nil
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
		f.league = append(f.league, Player{name,1})
	} else{
		player.Wins++
	}
	f.Database.Encode(&(f.league))
}

func (f *FileSystemStore) GetLeague() League {
	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Wins < f.league[j].Wins
	})
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

