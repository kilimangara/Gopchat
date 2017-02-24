package models

import "github.com/gorilla/websocket"

type Client struct {
	id int

	connection *websocket.Conn

	send chan []byte
}

type Message struct{
	author int

	message []byte

}

type Room struct {
	clients map[int]*Client

	inbound chan *Message

	register chan *Client

	unregister chan *Client
}

func NewRoom(ids []int)*Room{
	clients := make(map[int]*Client)
	for i:= range ids{
		append(clients, NewClient(i))
	}
	return &Room{
		clients:clients,
		inbound:make(chan *Message),
		register:make(chan *Client),
		unregister:make(chan *Client),
	}
}

func NewClient(id int) *Client{
	return &Client{
		id:id,
		connection:SearchConnection(id),
		send: make(chan []byte),
	}
}
