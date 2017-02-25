package handlers

import (
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"log"
	"Gopchat/models"
)

type DataJSON map[string]interface{}

var functions []func()

type msgFromClient struct {
	Type string `json:"type"`
	Data DataJSON `json:"data"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:1024,
	WriteBufferSize:1024,
	CheckOrigin:func(r *http.Request) bool {return true},
}

func ConnectionHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	userId, ok:= authenticate(r)

	if !ok{
		upgrader.Error(w, r, http.StatusBadRequest, "there is no such user with current token")
		return
	}

	conn, err:=upgrader.Upgrade(w, r, nil)

	if err!=nil{
		log.Fatal("handers.ConnectionHandler "+err.Error())
		return
	}

	client:=models.AttachConnectionToClient(userId, conn)


	client.StartClient()
}

func interceptInbondMessage(){

}


