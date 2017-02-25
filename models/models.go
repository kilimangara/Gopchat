package models

import (
	"github.com/gorilla/websocket"
	"Gopchat/db"
	"log"
)


type Client struct {
	Id int

	Connection *websocket.Conn

	Send chan []byte
}

type Message struct{
	Author int

	Message []byte

}

type Room struct {
	id int

	Clients map[int]*Client

	Inbound chan *Message

	Register chan *Client

	Unregister chan *Client

	Quit chan bool
}


func NewRoom(ids []int, id int)*Room{
	clients := make(map[int]*Client)
	for i:= range ids{
		append(clients, NewClient(i))
	}
	var room Room = Room{
		id:id,
		Clients:clients,
		Inbound:make(chan *Message),
		Register:make(chan *Client),
		Unregister:make(chan *Client),
		Quit:make(chan bool, 1),
	}
	go room.run()
	return &room
}

func NewClient(id int) *Client{
	return &Client{
		Id:id,
		Connection:nil,
		Send: make(chan []byte),
	}
}

func(client *Client) StartClient(){
	db.SwitchClientConnected(client.Id)
	defer db.SwitchClientDisconnected(client.Id)
	go client.writePump()
	client.readPump()

}

func(client *Client) writePump(){

}

func(client *Client) readPump(){

}

func (room *Room)run(){
	for{
		select {
		case client:=<-room.Register:
			room.Clients[client.Id]=client
		case client:=<-room.Unregister:
			delete(room.Clients, client.Id)
		case message:=<-room.Inbound:
			db.SaveMessage(message, room.id)
			for _,client:= range room.Clients{
				select {
				case client.Send <- message:
				}
			}
		case <-room.Quit:
			log.Println("room stopped "+room.id)
			return
		}
	}
}
