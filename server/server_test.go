package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)



func TestPlayerServer(t *testing.T) {
	store := &StubPlayerStore{
		map[string]string{
			"Tom":"20",
			"Nancy":"10",
		},
		nil,
	}
	server := &PlayerServer{store}
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

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	server := PlayerServer{NewInMemoryPlayerStore()}
	player := "Apollo"
	server.ServeHTTP(httptest.NewRecorder(),postNewRequest(player))
	server.ServeHTTP(httptest.NewRecorder(),postNewRequest(player))
	server.ServeHTTP(httptest.NewRecorder(),postNewRequest(player))
	response := httptest.NewRecorder()
	server.ServeHTTP(response,getNewRequest(player))
	assertResponseStatus(t,response.Code,http.StatusOK)
	assertString(t,response.Body.String(),"3")
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
