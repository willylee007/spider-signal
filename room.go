package main

import (
	"fmt"
)

// Room defined a room contain many clients
type Room struct {
	hub *Hub
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *ClientMsg

	relay chan *RelayClientSignal

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// room name
	name string

	id string
}

func newRoom(name, id string, hub *Hub) *Room {
	return &Room{
		name:       name,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *ClientMsg),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		relay:      make(chan *RelayClientSignal),
		id:         id,
		hub:        hub,
	}
}
func serveRoom(r *Room) {
	go broadcastMsg(r)
}
func broadcastMsg(r *Room) {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
			fmt.Printf("client %s enter room, room: %s ", client.id, r.id)
			for otherClient := range r.clients {
				if otherClient.id != client.id {
					queueMsg := otherClient.msgs

					for i := 0; i < len(queueMsg); i++ {
						client.send <- queueMsg[i]
					}
				}
			}
		case client := <-r.unregister:
			fmt.Printf("client %s leave room %s room name:%s\n", client.id, r.id, r.name)
			delete(r.clients, client)
			close(client.send)
			count := len(r.clients)
			if count == 0 {
				fmt.Printf("close room: %s room name:%s\n", r.id, r.name)
				r.hub.unregisterRoom <- r
			}
		case relayMsg := <-r.relay:
			srcClient := relayMsg.client
			msg := relayMsg.msg
			srcClient.msgs = append(srcClient.msgs, msg)

			for dstClient := range r.clients {
				if srcClient.id != dstClient.id {
					fmt.Println("转发消息")
					dstClient.send <- msg
				}
			}
		case <-r.broadcast:

		}
	}
}
