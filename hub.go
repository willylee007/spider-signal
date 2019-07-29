package main

import (
	"encoding/json"
	"fmt"

	"github.com/rs/xid"
)

const (
	joinMsg         = "join"
	leaveMsg        = "leave"
	chatMsg         = "chat"
	offerMsg        = "offer"
	answerMsg       = "answer"
	icecandidateMsg = "icecandidate"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// unregister room
	unregisterRoom chan *Room
	// Inbound messages from the clients.
	broadcast chan *ClientMsg

	unregisterClient chan *Client
	rooms            map[string]*Room
}

func newHub() *Hub {
	return &Hub{
		broadcast:        make(chan *ClientMsg),
		rooms:            make(map[string]*Room),
		unregisterRoom:   make(chan *Room),
		unregisterClient: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.unregisterClient:
			if client.room != nil {
				client.room.unregister <- client
			}
		case room := <-h.unregisterRoom:
			delete(h.rooms, room.name)
			fmt.Printf("clean room, hub have %d room currently \n", len(h.rooms))
		case message := <-h.broadcast:
			signal := SignalMsg{}
			client := message.client
			msg := message.msg
			err := json.Unmarshal(msg, &signal)
			if err != nil {
				return
			}

			switch signal.Type {
			case joinMsg:
				roomName := signal.RoomName
				var room = h.rooms[roomName]
				if room == nil {
					room = newRoom(roomName, xid.New().String(), h)
					h.rooms[roomName] = room
					client.room = room
					serveRoom(room)
				} else {
					client.room = room
				}
				room.register <- client
			case leaveMsg:
			case icecandidateMsg:
				fmt.Println("receive ice candidate")
				relayMsg := &RelayClientSignal{
					client: client,
					msg:    msg,
				}
				client.room.relay <- relayMsg
			case answerMsg:
				fmt.Println("receive answer")
				relayMsg := &RelayClientSignal{
					client: client,
					msg:    msg,
				}
				client.room.relay <- relayMsg
			case offerMsg:
				fmt.Println("receive offer")
				relayMsg := &RelayClientSignal{
					client: client,
					msg:    msg,
				}
				client.room.relay <- relayMsg
			case chatMsg:
				chatMsg := signal.Msg
				chatByte := []byte(chatMsg)
				queueMsg := &ClientMsg{
					client: client,
					msg:    chatByte,
				}
				client.room.broadcast <- queueMsg
			}

		}
	}
}
