package models

import (
	"github.com/gorilla/websocket"
	"Gopchat/db"
	"log"
)

var(
	roomPool map[int]*Room
	connectionsPools map[int]*websocket.Conn
)

func GetRoomById(id int)(*Room, error){
	var room *Room
	if room = roomPool[id]; room==nil{
		ids, err := db.GetRoom(id)
		if err!=nil{
			log.Fatal("pool.GetRoomById "+err)
			return nil, err
		}
		room = NewRoom(ids)
	}
	return room, nil
}

func SearchConnection(id int)*websocket.Conn{
	return connectionsPools[id]
}

