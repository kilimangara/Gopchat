package models

import (
	"Gopchat/db"
	"log"
	"sync"
	"github.com/gorilla/websocket"
)

var(
	roomPool map[int]*Room = make(map[int]*Room)
	connectionsPools map[int]*Client = make(map[int]*Client)
	connectionMutex *sync.RWMutex = sync.RWMutex{}
)

func GetRoomById(id int)(*Room, error){
	connectionMutex.RLock()
	defer connectionMutex.RUnlock()
	var room *Room
	if room = roomPool[id]; room==nil{
		ids, err := db.GetRoom(id)
		if err!=nil{
			log.Fatal("pool.GetRoomById "+err)
			return nil, err
		}
		room = NewRoom(ids,id)
	}
	return room, nil
}

func SearchConnection(id int)*Client{
	connectionMutex.RLock()
	defer connectionMutex.RUnlock()
	if(connectionsPools[id]==nil){
		connectionsPools[id]= NewClient(id)
	}
	return connectionsPools[id]
}

func AttachConnectionToClient(id int, conn *websocket.Conn)(*Client){
	client:=SearchConnection(id)
	client.Connection = conn
	return client
}

func DetachConnectionFromClient(id int)(*Client){
	client:=SearchConnection(id)
	client.Connection=nil
	return client
}

