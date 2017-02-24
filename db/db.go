package db

import (
	"database/sql"
	"Gopchat/builder"
	"log"
	"Gopchat/models"
	"time"
)

var(
	db *sql.DB
)

const(
	CLIENTS_ROOM_BY_ID="SELECT user_room.id_user FROM user_room WHERE id_room=$1"
	SAVE_MESSAGE="INSERT INTO messages VALUES(default, $1, $2,$3, $4)"
)


func Init(){
	db, _ = builder.BuildPostgresClient()
	if err:=db.Ping(); err!=nil{
		log.Fatal(err)
	}
}

func GetRoom(id int)([]int, error){
	result := make([]int, 10)
	rows, err := db.Query(CLIENTS_ROOM_BY_ID, id)
	defer rows.Close()
	if err!=nil{
		log.Fatal("db.GetRoom "+err)
		return result, err
	}
	for rows.Next(){
		var i int
		if err:=rows.Scan(&i); err!=nil{
			log.Fatal("db.GetRoom "+err)
			return result, err
		}
		append(result, i)
	}
	return result, nil
}

func SaveMessage(message *models.Message, idRoom int) error{
	_, err:=db.Exec(SAVE_MESSAGE, &message.Author, idRoom, string(&message.Message), time.Now().Format(time.RFC3339))
	if(err!=nil){
		log.Fatal("db.SaveMessage "+err)
		return err
	}
	return nil
}

func Close(){
	db.Close()
}
