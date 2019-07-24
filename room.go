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
	broadcast chan []byte

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
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
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
			fmt.Printf("client %s enter room, room: %s room name:%s, room register number:%d \n ", client.id, r.id, r.name, len(r.clients))
		case client := <-r.unregister:
			fmt.Printf("client %s leave room %s room name:%s\n", client.id, r.id, r.name)
			delete(r.clients, client)
			close(client.send)
			count := len(r.clients)
			if count == 0 {
				fmt.Printf("close room: %s room name:%s\n", r.id, r.name)
				r.hub.unregisterRoom <- r
			}
		case message := <-r.broadcast:
			for client := range r.clients {
				client.send <- message
			}
		}
	}
}
