package server

type StubPlayerStore struct {
	Scores map[string] string
	WinCalls []string
	league []Player
}
func (sp StubPlayerStore)GetPlayerScore(player string) string{
	return sp.Scores[player]
}
func (sp *StubPlayerStore)RecordWin(name string) {
	sp.WinCalls = append(sp.WinCalls, name)
}
func (sp *StubPlayerStore)GetLeague() League {
	return sp.league
}