package elo

type Club struct {
	matches []*Match
	players *Players
	name string

}

func (club *Club) addMatch(match *Match){
	
}

func (club *Club) addPlayer(name string) bool {
	p := club.players.New(name)
	result := club.players.AddPlayer(p)
	return result
}