package models

import (
	"appengine/datastore"
	"common"
)

const (
	GameKind    = "Game"
	AllGamesKey = "Games{All}"
)

type State string

const (
	StateCreated  = "Created"
	StatePlaying  = "Playing"
	StateFinished = "Finished"
)

type Games []Game

func (self Games) process(c common.Context) Games {
	for index, _ := range self {
		(&self[index]).process(c)
	}
	return self
}

type Game struct {
	Id          *datastore.Key
	Players     []*datastore.Key
	State       State
	PlayerNames []string `datastore:"-"`
}

func (self *Game) process(c common.Context) *Game {
	self.PlayerNames = make([]string, len(self.Players))
	for index, id := range self.Players {
		if ai := GetAIById(c, id); ai != nil {
			self.PlayerNames[index] = ai.Name
		} else {
			self.PlayerNames[index] = "[redacted]"
		}
	}
	return self
}

func findAllGames(c common.Context) (result Games) {
	ids, err := datastore.NewQuery(GameKind).GetAll(c, &result)
	common.AssertOkError(err)
	for index, id := range ids {
		result[index].Id = id
	}
	if result == nil {
		result = Games{}
	}
	return
}

func GetAllGames(c common.Context) (result Games) {
	common.Memoize(c, AllGamesKey, &result, func() interface{} {
		return findAllGames(c)
	})
	return result.process(c)
}

func (self *Game) Save(c common.Context) *Game {
	var err error
	if self.Id == nil {
		self.State = StateCreated
		self.Id, err = datastore.Put(c, datastore.NewKey(c, GameKind, "", 0, nil), self)
	} else {
		_, err = datastore.Put(c, self.Id, self)
	}
	common.AssertOkError(err)
	common.MemDel(c, AllGamesKey)
	return self
}
