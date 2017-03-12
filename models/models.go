package models

import (
	"github.com/gorilla/websocket"
	"Gopchat/db"
	"log"
	"Gopchat/handlers"
	"time"
	"bytes"
	"encoding/json"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512

	NEW_MESSAGE="NEW_MESSAGE"

	NEW_USER ="NEW_USER"

)
var(
	newLine = []byte{'\n'}
	space = []byte{' '}
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
	var room= Room{
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
		Connection:searchConnection(id),
		Send: make(chan []byte),
	}
}

func NewEmptyClient(id int)*Client{
	return &Client{
		Id:id,
		Connection:nil,
		Send:make(chan []byte),
	}
}

func(client *Client) StartClient(interceptors []func(id int, msg *handlers.MsgFromClient)(bool, error)){
	db.SwitchClientConnected(client.Id)
	defer db.SwitchClientDisconnected(client.Id)
	go client.writePump()
	client.readPump(interceptors)

}

func buildResponse(typeOfEvent string)([]byte, error){
	var response map[string]string = make(map[string]string)
	response["type"]=typeOfEvent
	return json.Marshal(response)
}

func(client *Client) writePump(){

	ticker:=time.NewTicker(pingPeriod)
	 defer func(){
		 ticker.Stop()
		 DetachConnectionFromClient(client.Id)
	 }()
	for {
		select{
		case message, ok:= <- client.Send:
			client.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok{
				client.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err:= client.Connection.NextWriter(websocket.TextMessage)
			if err!=nil{
				return
			}
			w.Write(message)
		}
	}
}

func(client *Client) readPump(interceptors []func(id int, msg *handlers.MsgFromClient)(bool, error)){
	defer DetachConnectionFromClient(client.Id)
	client.Connection.SetReadLimit(maxMessageSize)
	client.Connection.SetReadDeadline(time.Now().Add(pongWait))
	client.Connection.SetPongHandler(func(data string) error { client.Connection.SetReadDeadline(time.Now().Add(pongWait)); return nil})
	var msgStruct *handlers.MsgFromClient
	for {
		_,message,err:=client.Connection.ReadMessage()
		if(err!=nil){
			if(websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway)) {
				log.Println(err.Error())
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message,newLine, space,-1))
		if errr:= json.Unmarshal(message, &msgStruct); errr!=nil{
			errorMsg:= handlers.NewMsg(handlers.ERROR_BAD_FORMAT,"can't convert your shit")
			client.Send<-[]bytes(errorMsg.Error())
			continue
		}
		for _,interceptor:= range interceptors{
			flag,err:= interceptor(client.Id, msgStruct)
			if(flag) {
				if (err != nil) {
					client.Send <- []bytes(err.Error())
				}
				break
			}
		}
	}


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
			msg,_:=buildResponse(NEW_MESSAGE)
			for _,client:= range room.Clients{
				if(client.Id != message.Author){
					client<-msg
				}
			}
		case <-room.Quit:
			log.Println("room stopped "+room.id)
			return
		}
	}
}
