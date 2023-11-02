package main

import (
	"log"
	"sync"
)

type Hub struct {
	Clients    map[*Client]ClientAttributes
	Register   chan *Client
	Unrigester chan *Client
	Broadcast  chan *Recipent
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]ClientAttributes),
		Register:   make(chan *Client),
		Unrigester: make(chan *Client),
		Broadcast:  make(chan *Recipent),
	}
}

func (h *Hub) Run() {
	var mutex sync.Mutex
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = ClientAttributes{Hubb: h, Conn: client.ClientAttributes.Conn}
			name := client.Name
			log.Printf("Client %v has joined the chat", name)
			for client := range h.Clients {
				userDetail := &RegAndUnregForm{
					Username: name,
					Message:  name + " is a new member of the chat",
				}

				for i, j := range h.Clients {
					k := ""
					k += i.Name
					log.Println(k, j)
				}
				h.Clients[client].Conn.WriteJSON(userDetail)
			}
		case client := <-h.Unrigester:
			name := client.Name
			delete(h.Clients, client)
			log.Printf("Client %v has left the chat", name)
			for client := range h.Clients {
				userDetail := &RegAndUnregForm{
					Username: name,
					Message:  name + " has left the chat",
				}
				h.Clients[client].Conn.WriteJSON(userDetail)
			}
		case message := <-h.Broadcast:

			flag := false
			for client := range h.Clients {
				if message.Username == client.Name {
					mutex.Lock()
					h.Clients[client].Conn.WriteJSON(message)
					mutex.Unlock()
					flag = true

					break
				}

			}
			if message.Username == "Everyone" {
				flag = true
				for client := range h.Clients {

					mutex.Lock()
					h.Clients[client].Conn.WriteJSON(message)
					mutex.Unlock()
				}
				break
			}

			if !flag {
				for client := range h.Clients {
					if client.Name == message.Sender && message.Username != "Everyone" {
						Response := &NoUserFound{
							Status: 404,
							StatusCode:  "UserNotFound",
							StatusMessage: "The user you are trying to connect is no longer available or never part of the chat",
						}
						mutex.Lock()
						h.Clients[client].Conn.WriteJSON(Response)
						mutex.Unlock()
					}
				}
			}
		}

	}
}
