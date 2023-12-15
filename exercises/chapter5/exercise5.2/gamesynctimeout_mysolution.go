package main

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"
)

type Player struct {
	Id int
}

type GameLoader struct {
	QueuedPLayers []*Player
	PlayersReady  map[*Player]bool
	LoadReadyCond *sync.Cond
	LoadTimedOut  bool
}

func NewGameLoader(players []*Player) *GameLoader {
	self := GameLoader{
		QueuedPLayers: players,
		PlayersReady:  make(map[*Player]bool),
		LoadReadyCond: sync.NewCond(&sync.Mutex{}),
		LoadTimedOut:  false,
	}
	for _, player := range players {
		self.PlayersReady[player] = false
	}
	return &self
}

func (x *GameLoader) areAllPlayersReady() bool {
	for _, ready := range x.PlayersReady {
		if !ready {
			return false
		}
	}
	return true
}

func (x *GameLoader) ScheduleTimeout(duration time.Duration) {
	fmt.Println("Scheduling timeout...")
	time.Sleep(duration)

	x.LoadReadyCond.L.Lock()

	if !x.areAllPlayersReady() {
		fmt.Println("Timed out!")
		x.LoadReadyCond.Broadcast()
		x.LoadTimedOut = true
	}

	x.LoadReadyCond.L.Unlock()
}

func (x *GameLoader) WirePlayerIn(player *Player) error {
	if !slices.Contains(x.QueuedPLayers, player) {
		return errors.New("invalid player")
	}
	if x.PlayersReady[player] {
		return errors.New(fmt.Sprintf("player [%d] is already ready", player.Id))
	}
	if x.LoadTimedOut {
		fmt.Println(player.Id, ": Aborted.")
		return nil
	}

	x.LoadReadyCond.L.Lock()

	x.PlayersReady[player] = true
	fmt.Println(player.Id, ": Connected")

	if x.areAllPlayersReady() {
		// fmt.Println("All players connected. Ready player", player.Id)
		x.LoadReadyCond.Broadcast()
	}
	for !x.areAllPlayersReady() && !x.LoadTimedOut {
		fmt.Println(player.Id, ": Waiting for other players")
		x.LoadReadyCond.Wait()
	}

	x.LoadReadyCond.L.Unlock()

	/* without the below lines of code, stdout output would be like:
	```
	0 : Connected
	0 : Waiting for other players
	Scheduling timeout...
	1 : Connected
	1 : Waiting for other players
	Timed out!
	1 : Waiting for other players
	0 : Waiting for other players
	2 : Connected
	2 : Waiting for other players
	3 : Connected
	All players connected. Ready player 3
	```
	*/

	// because cond waits for 2 conditions
	if x.LoadTimedOut {
		fmt.Println(player.Id, ": Game cancelled.")
	} else {
		fmt.Println("All players connected. Ready player", player.Id)
	}
	return nil
}

func main() {
	playersInGame := 4

	var players []*Player
	for playerId := 0; playerId < playersInGame; playerId++ {
		players = append(players, &Player{Id: playerId})
	}

	gameState := NewGameLoader(players)

	go gameState.ScheduleTimeout(time.Second)
	for _, player := range players {
		go func() {
			var err = gameState.WirePlayerIn(player)
			if err != nil {
			}
		}()
		time.Sleep(800 * time.Millisecond)
	}
	time.Sleep(time.Minute)
}
