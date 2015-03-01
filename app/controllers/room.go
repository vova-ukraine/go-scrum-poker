package controllers

import (
	"code.google.com/p/go.net/websocket"
	"github.com/revel/revel"
	"math/rand"
	"strconv"
	"fmt"
	"scrum-poker/app/libs"
)

var (
	rooms = make(map[string] libs.PokerRoom)
)

type Room struct {
	*revel.Controller
}

func (c Room) Create() revel.Result {

	var id = strconv.Itoa(rand.Int())
	_, ok := rooms[id]

	// If room exists lets try again
	if (ok) { return c.Redirect(App.Index) }

	room := libs.PokerRoom{Admin: c.Session.Id()}
	room.Members = make(map[string]libs.PokerMember)
	rooms[id] = room

	return c.Redirect("/room/%s", id)
}

func (c Room) Delete(id string) revel.Result  {
	_, ok := rooms[id]

	if (ok) {
		delete(rooms, id)
		//@TODO: Notify all subscribers
	}
	return c.Redirect(App.Index)
}

func (c Room) Index(id string) revel.Result {
	room, ok := rooms[id]
	if (!ok) { return c.Redirect(App.Index) }
	var admin = (room.Admin == c.Session.Id())
	return c.Render(id, admin)
}

func (c Room) Socket(id string, ws *websocket.Conn) revel.Result {
fmt.Println("")
	subscribtion := make(chan string)
	member := libs.PokerMember{"", subscribtion}
	rooms[id].Members[c.Session.Id()] = member

	defer rooms[id].DeleteMember(c.Session.Id())

	votes := make(chan string)
	go func() {
		var msg string
		for {
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				close(votes)
				return
			}
			votes <- msg
		}
	}()

	go rooms[id].NotifyMembers()

	for {
		select {
		case vote, ok := <-votes:
			// If the channel is closed, they disconnected.
			if !ok {
				return nil
			}
			rooms[id].SetVote(c.Session.Id(), vote)
		case <-subscribtion:

			if websocket.JSON.Send(ws, rooms[id].GetJSONState(c.Session.Id())) != nil {
				// They disconnected.
				return nil
			}
		}
	}
	return nil
}
