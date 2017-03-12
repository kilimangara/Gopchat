package models

import (
	"Gopchat/db"
	"log"
	"sync"
	"github.com/gorilla/websocket"
	"Gopchat/tokens"
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
		roomPool[id] = room
	}
	return room, nil
}

func searchConnection(id int)*Client{
	if(connectionsPools[id]==nil){
		connectionsPools[id]= NewEmptyClient(id)
	}
	return connectionsPools[id]
}

func AttachConnectionToClient(id int, conn *websocket.Conn)(*Client){
	connectionMutex.Lock()
	defer connectionMutex.Unlock()
	client:=searchConnection(id)
	client.Connection = conn
	return client
}

func DetachConnectionFromClient(id int)(*Client){
	connectionMutex.Lock()
	defer connectionMutex.Unlock()
	client:=searchConnection(id)
	client.Connection.Close()
	client.Connection=nil
	return client
}

func AddRoomToPool(roomPattern tokens.JSONMsg){
	var clients []int =roomPattern["clients"]
	var id int = roomPattern["id"]
	connectionMutex.Lock()
	defer connectionMutex.Unlock()
	room:=NewRoom(clients, id)
	roomPool[room.id]=room
}

func AddClientToRoom(userPattern *tokens.JSONMsg){
	//adding user to some room
}

