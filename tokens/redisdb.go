package tokens

import (
	"gopkg.in/redis.v5"
	"Gopchat/builder"
	"log"
)

var(
	db *redis.Client
)

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
		msg, err1:= pubsub.ReceiveMessage()
		if(err1!=nil){
			//some error
			continue
		}

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
