package server



type StubPlayerStore struct {
	Scores map[string] string
	WinCalls []string
}
func (sp StubPlayerStore)GetPlayerScore(player string) string{
	return sp.Scores[player]
}
func (sp *StubPlayerStore)RecordWin(name string) {
	sp.WinCalls = append(sp.WinCalls, name)
}
