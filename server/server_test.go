package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

var ApplicationJson = "application-json"
var ContentType = "content-type"
func TestPlayerServer(t *testing.T) {
	store := &StubPlayerStore{
		map[string]string{
			"Tom":"20",
			"Nancy":"10",
		},
		nil,
		League{},
	}
	server := NewPlayerServer(store)
	t.Run("返回Tom的游戏分数", func(t *testing.T) {
		player := "Tom"
		request:= getNewRequest(player)
		response := httptest.NewRecorder()
		server.ServeHTTP(response,request)
		assertResponseStatus(t,response.Code,http.StatusOK)
		assertString(t,response.Body.String(),"20")
	})
	t.Run("返回Nancy的游戏分数", func(t *testing.T) {
		player := "Nancy"
		request:= getNewRequest(player)
		response := httptest.NewRecorder()
		server.ServeHTTP(response,request)
		assertResponseStatus(t,response.Code,http.StatusOK)
		assertString(t,response.Body.String(),"10")
	})
	t.Run("用户不存在", func(t *testing.T) {
		player := "Apollo"
		request:= getNewRequest(player)
		response := httptest.NewRecorder()
		server.ServeHTTP(response,request)
		assertResponseStatus(t,response.Code,http.StatusNotFound)
	})
	t.Run("当收到POST请求时,记录获胜记录", func(t *testing.T) {
		player := "Apollo"
		request:= postNewRequest(player)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertResponseStatus(t, response.Code, http.StatusAccepted)
		if len(store.WinCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d",len(store.WinCalls),1)
		}
		if store.WinCalls[0] != "Apollo" {
			t.Errorf("got %s,want Apollo",store.WinCalls[0])
		}
	})
}

func TestLeague(t *testing.T) {
	store := &StubPlayerStore{
		map[string]string{
			"Tom":"20",
			"Nancy":"10",
		},
		nil,
		[]Player{{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},},
	}
	server := NewPlayerServer(store)
	t.Run("对于/league路由，服务器会返回一个200状态码", func(t *testing.T) {
		request,_ := http.NewRequest(http.MethodGet,"/league", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response,request)
		var got []Player
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Errorf("response: %s, err: %#v",response.Body,err)
		}
		assertResponseStatus(t,response.Code,http.StatusOK)
	})
	t.Run("对于/league路由，服务器会返回Json格式的玩家列表",func (t *testing.T) {
		request,_ := http.NewRequest(http.MethodGet,"/league", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response,request)
		want := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}
		var got []Player
		err := json.NewDecoder(response.Body).Decode(&got)
		assertErrShouldBeNil(t,err)
		assertResponseHeader(t,response,ContentType,ApplicationJson)
		assertResponseStatus(t,response.Code,http.StatusOK)
		assertJson(t,got,want)
	})
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database,cleanDatabase := createTempfile(t,`[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)
	defer cleanDatabase()
	store,err:= NewFileSystemStore(database)
	assertErrShouldBeNil(t,err)
	server:= NewPlayerServer(store)
	player := "Apollo"
	postPlayerServerTimes(server,player,3)
	response := httptest.NewRecorder()
	server.ServeHTTP(response,getNewRequest(player))
	assertResponseStatus(t,response.Code,http.StatusOK)
	assertString(t,response.Body.String(),"3")
}

func TestRecordingWinsAndRetrievingLeague(t *testing.T) {
	database,cleanDatabase := createTempfile(t,"[]")
	defer cleanDatabase()
	store,err:= NewFileSystemStore(database)
	assertErrShouldBeNil(t,err)
	server:= NewPlayerServer(store)
	player := "Tom"
	postPlayerServerTimes(server,player,3)
	response := httptest.NewRecorder()
	request,err := http.NewRequest(http.MethodGet,"/league", nil)
	server.ServeHTTP(response,request)
	var got []Player
	want := []Player{
		{"Tom",3},
	}
	err = json.NewDecoder(response.Body).Decode(&got)
	assertErrShouldBeNil(t,err)
	assertResponseHeader(t,response,ContentType,ApplicationJson)
	assertResponseStatus(t,response.Code,http.StatusOK)
	assertJson(t,got,want)
}

func TestFileSystemStore(t *testing.T) {
	t.Run("通过文件Reader获取/league的返回", func(t *testing.T) {
		database,cleanDatabase := createTempfile(t,`[
            {"Name": "Cleo", "Wins": 33},
            {"Name": "Chris", "Wins": 10}]`)
		defer cleanDatabase()
		store ,_:= NewFileSystemStore(database)
		got := store.GetLeague()
		want := []Player{
			{"Chris",10},
			{"Cleo",33},
		}
		assertLeague(t,got,want)
		// read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})
	t.Run("通过文件Reader获取 getPlayerScores", func(t *testing.T) {
		database,cleanDatabase := createTempfile(t,`[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()
		store ,_:= NewFileSystemStore(database)
		got := store.GetPlayerScore("Cleo")
		want := "10"
		assertString(t,got, want)
	})
	t.Run("为已有玩家添加一次胜利", func(t *testing.T) {
		database,cleanDatabase := createTempfile(t,`[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()
		store ,_:= NewFileSystemStore(database)
		store.RecordWin("Chris")
		got := store.GetPlayerScore("Chris")
		want := "34"
		assertString(t,got, want)
	})
	t.Run("为新玩家添加胜利", func(t *testing.T) {
		database,cleanDatabase := createTempfile(t,`[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()
		store ,_:= NewFileSystemStore(database)
		store.RecordWin("Tom")
		got := store.GetPlayerScore("Tom")
		want := "1"
		assertString(t,got, want)
	})
	t.Run("处理空文件", func(t *testing.T) {
		database,cleanDatabase := createTempfile(t,"")
		defer cleanDatabase()
		_,err:= NewFileSystemStore(database)
		assertErrShouldBeNil(t,err)
	})
}

func createTempfile(t *testing.T,initialData string) (*os.File,func()) {
	t.Helper()
	tmpfile,err := ioutil.TempFile("","db")
	if err != nil {
		t.Errorf("创建临时文件失败:%s",err.Error())
	}
	tmpfile.Write([]byte(initialData))
	removeFile := func() {
		os.Remove(tmpfile.Name())
	}
	return tmpfile,removeFile
}

func assertLeague(t *testing.T, got []Player, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %#v,want %#v",got,want)
	}
}

func postPlayerServerTimes(server *PlayerServer,player string,times int) {
	for i := 0;i < times;i++ {
		server.ServeHTTP(httptest.NewRecorder(),postNewRequest(player))
	}
}

func getNewRequest(player string) *http.Request {
	request,_ := http.NewRequest(http.MethodGet,fmt.Sprintf("/players/%s",player),nil)
	return request
}

func postNewRequest(player string) *http.Request {
	request,_ := http.NewRequest(http.MethodPost,fmt.Sprintf("/players/%s",player),nil)
	return request
}

func assertResponseHeader(t *testing.T,response *httptest.ResponseRecorder,key string,want string) {
	if response.Header().Get(key) != want {
		t.Errorf("response did not have content-type of application/json, got %v", response.Header())
	}
}

func assertJson(t *testing.T, got []Player, want []Player) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v , want %#v",got,want)
	}
}

func assertErrShouldBeNil(t *testing.T,err error) {
	if err != nil {
		t.Errorf("err: %#v",err)
	}
}

func assertString(t *testing.T,got string,want string){
	t.Helper()
	if got!= want {
		t.Errorf("Got %s, Want %s",got, want)
	}
}

func assertResponseStatus(t *testing.T,got int,want int) {
	t.Helper()
	if got!= want {
		t.Errorf("got %#v,want %#v",got,want)
	}
}

