package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	Name             string
	ClientAttributes *ClientAttributes
}

type ClientAttributes struct {
	Hubb *Hub
	Conn *websocket.Conn
}

type RegAndUnregForm struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type Recipent struct {
	Sender   string `json:"Sender"`
	Username string `json:"Recipient"`
	Message  string `json:"Message"`
}
type NoUserFound struct {
	//UserNotFound string `json:"UserNotFound"`
	Status       int    `json:"Status"`
	StatusCode string `json:"StatusCode"`
	StatusMessage string `json:"StatusMessage"`
	//ErrorCode map[int]string `json:NotExistingUser"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *Client) recoveryDetails() {
	c.ClientAttributes.Hubb.Unrigester <- c
	c.ClientAttributes.Conn.Close()
	if r := recover(); r != nil {
		log.Println("Still panicing, err:", r)
	}
}

func (c *Client) Read() {
	defer c.recoveryDetails()
	for {
		var msg Recipent
		_, message, err := c.ClientAttributes.Conn.ReadMessage()
		if err != nil {
			log.Panic("Message is not read, err:", err)
		}

		log.Println("[message]: ", string(message))
		json.Unmarshal(message, &msg)
		c.ClientAttributes.Hubb.Broadcast <- &msg

		//c.ClientAttributes.Conn.WriteMessage(websocket.TextMessage, message)

	}

}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("name")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panic("Connection not upgraded, err:", err)
	}
	client := &Client{Name: userName}
	client.ClientAttributes = &ClientAttributes{Hubb: hub, Conn: conn}
	hub.Register <- client

	client.Read()
}
