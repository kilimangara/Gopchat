package tokens

import (
	"gopkg.in/redis.v5"
	"Gopchat/builder"
	"log"
	"Gopchat/models"
	"encoding/json"
)

var(
	db *redis.Client
)

type JSONMsg map[string]interface{}

const(
	ROOM_CHANNEL="ROOM"
	USER_CHANNEL="USER"
)

func Init(){
	db,_= builder.BuildRedisClient()
	if err := db.Ping().Err(); err!=nil{
		log.Fatal(err.Error())
	}
}

func SubscribeToRoomChannel(){
	pubsub, err:= db.Subscribe(ROOM_CHANNEL)
	if err!=nil{
		log.Fatal("pizda "+err.Error())
	}
	for{
		msg, err:= pubsub.ReceiveMessage()
		if(err!=nil){
			//some error
			continue
		}
		room:= &JSONMsg{
			"id":0,
			"clients":[]int{},
		}
		err2:=json.Unmarshal([]byte(msg), &room)
		if(err2!=nil){
			//some error
			continue
		}
		models.AddRoomToPool(room)
	}
}

func SubscribeToUserChannel(){

}


func GetUser(token string)(int, bool){
	userId, err:= db.Get(token).Int64()

	if err !=nil{
		if err!=redis.Nil{
			log.Fatal(err.Error())
		}
		return 0, false
	}

	return userId, true
}

func Close(){
	db.Close()
}
