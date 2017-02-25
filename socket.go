package main

import (
	"Gopchat/db"
	"Gopchat/tokens"
	"github.com/julienschmidt/httprouter"
	"Gopchat/handlers"
	"log"
	"net/http"
)

func main(){
	db.Init()
	defer db.Close()
	tokens.Init()
	defer tokens.Close()
	router:=httprouter.New()
	router.GET("/ws", handlers.ConnectionHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}
