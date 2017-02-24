package db

import (
	"database/sql"
	"Gopchat/builder"
	"log"
)

var(
	db *sql.DB
)

const(
	ROOM_BY_ID="SELECT user_room.id_user FROM user_room WHERE id_room=$1"
)


func Init(){
	db, _ = builder.BuildPostgresClient()
	if err:=db.Ping(); err!=nil{
		log.Fatal(err)
	}
}

func GetRoom(id int)([]int, error){
	result := make([]int, 10)
	rows, err := db.Query(ROOM_BY_ID, id)
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

func Close(){
	db.Close()
}
