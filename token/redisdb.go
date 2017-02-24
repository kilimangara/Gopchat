package token

import (
	"gopkg.in/redis.v5"
	"Gopchat/builder"
	"log"
)

var(
	db *redis.Client
)

func Init(){
	db,_= builder.BuildRedisClient()
	if err := db.Ping().Err(); err!=nil{
		log.Fatal(err.Error())
	}
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
