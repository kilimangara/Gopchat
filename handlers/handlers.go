package handlers

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"log"
	"Gopchat/models"
	"github.com/gorilla/websocket"
	"encoding/json"
	"fmt"
)
const(
	SEND_MESSAGE="SEND_MESSAGE"
	ERROR_NO_SUCH_ROOM="NO_SUCH_ROOM"
	ERROR_BAD_FORMAT="BAD_FORMAT"

)

type DataJSON map[string]interface{}

//var functions []func(id int, msg *MsgFromClient)(bool, error)

type MsgFromClient struct {
	Type string `json:"type"`
	Data DataJSON `json:"data"`
}

func NewMsg(typeError, description string)(*MsgFromClient){
	return &MsgFromClient{
		Type:typeError,
		Data:DataJSON{
			"description":description,
		},
	}
}

func(msg *MsgFromClient)Error() string{
	strMsg,_:= json.Marshal(msg)
	return string(strMsg)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:1024,
	WriteBufferSize:1024,
	CheckOrigin:func(r *http.Request) bool {return true},
	Error:func(w http.ResponseWriter, r *http.Request, status int, reason string) {
		w.Header().Set("Sec-Websocket-Version", "13")
		http.Error(w, reason, status)
	},
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


	client.StartClient([]func(id int, msg *MsgFromClient)(bool,error){interceptInbondMessage})
}

func createErrorResponse(typeError string)(*MsgFromClient){
	var msg MsgFromClient = MsgFromClient{
		Type:typeError,
	}
	return &msg
}

func interceptInbondMessage(id int,msg *MsgFromClient)(bool,error){
	if msg.Type==SEND_MESSAGE {
		var to = msg.Data["to"]
		room, err:=models.GetRoomById(to)
		if err!=nil{
			return true, NewMsg(ERROR_NO_SUCH_ROOM,fmt.Sprintf("there is no room with id=%d",to))
		}
		byteMessage,err1 := json.Marshal(msg.Data["message"])
		if err1!=nil{
			return true, NewMsg(ERROR_BAD_FORMAT, "there is no message field in json Data")
		}
		message:= models.Message{
			Author:id,
			Message:byteMessage,
		}
		room.Inbound<-&message
		return true, nil
	}
	return false,nil
}

func interceptSmt(id int, msg *MsgFromClient)(bool, error){
	return true, nil
}


