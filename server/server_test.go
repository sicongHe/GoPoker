package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)



func TestPlayerServer(t *testing.T) {
	store := &StubPlayerStore{
		map[string]string{
			"Tom":"20",
			"Nancy":"10",
		},
		nil,
		[]Player{},
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
		if err != nil {
			t.Errorf("response: %s, err: %#v",response.Body,err)
		}
		if response.Header().Get("content-type") != "application-json" {
			t.Errorf("response did not have content-type of application/json, got %v", response.Header())
		}
		assertResponseStatus(t,response.Code,http.StatusOK)
		assertJson(t,got,want)
	})
}


func assertJson(t *testing.T, got []Player, want []Player) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v , want %#v",got,want)
	}
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	server := NewPlayerServer(NewInMemoryPlayerStore())
	player := "Apollo"
	server.ServeHTTP(httptest.NewRecorder(),postNewRequest(player))
	server.ServeHTTP(httptest.NewRecorder(),postNewRequest(player))
	server.ServeHTTP(httptest.NewRecorder(),postNewRequest(player))
	response := httptest.NewRecorder()
	server.ServeHTTP(response,getNewRequest(player))
	assertResponseStatus(t,response.Code,http.StatusOK)
	assertString(t,response.Body.String(),"3")
}

func TestRecordingWinsAndRetrievingLeague(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server:= NewPlayerServer(store)
	player := "Tom"
	server.ServeHTTP(httptest.NewRecorder(),postNewRequest(player))
	server.ServeHTTP(httptest.NewRecorder(),postNewRequest(player))
	server.ServeHTTP(httptest.NewRecorder(),postNewRequest(player))
	response := httptest.NewRecorder()
	request,_ := http.NewRequest(http.MethodGet,"/league", nil)
	server.ServeHTTP(response,request)
	var got []Player
	want := []Player{
		{"Tom",3},
	}
	err := json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Errorf("response: %s, err: %#v",response.Body,err)
	}
	if response.Header().Get("content-type") != "application-json" {
		t.Errorf("response did not have content-type of application/json, got %v", response.Header())
	}
	assertResponseStatus(t,response.Code,http.StatusOK)
	assertJson(t,got,want)
}

func getNewRequest(player string) *http.Request {
	request,_ := http.NewRequest(http.MethodGet,fmt.Sprintf("/players/%s",player),nil)
	return request
}

func postNewRequest(player string) *http.Request {
	request,_ := http.NewRequest(http.MethodPost,fmt.Sprintf("/players/%s",player),nil)
	return request
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

