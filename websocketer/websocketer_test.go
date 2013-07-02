package websocketer

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"testing"	
	"time"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
)

func TestText(t *testing.T) {
	//setup test variables
	var isOpened,isClosed bool
	message := ""
	
	//create websocket handler
	ws := New()
	ws.OnOpen(rack.Func(func(vars map[string]interface{},next func()) {
		isOpened = true
		next()
	}))
	ws.OnMessage(rack.Func(func(vars map[string]interface{},next func()) {
		message = V(vars).GetMessage().(string)
		V(vars).SetResponse("Who's there?")
		next()
	}))
	ws.OnClose(rack.Func(func(vars map[string]interface{},next func()) {
		isClosed = true
		next()
	}))
	ws.ReceiveTextMessages()
	
	//start a web server
	server := httper.HttpConnection(":4016")
	go server.Go(ws)
	
	//connect via websocket
	conn,err := websocket.Dial("ws://localhost:4016","","http://localhost")
	if err != nil {
		t.Fatal("Fatal - Can't open")
	}
	
	if !isOpened {
		t.Error("Should be open")
	}
	
	//send a message
	if _,err := conn.Write([]byte("Knock knock!")); err != nil {
		t.Fatal("Fatal - Can't write")
	}
	
	//get a response
	response := make([]byte,512)
	n, err := conn.Read(response)
	if err != nil {
		t.Fatal("Fatal - Can't read")
	}
	
	if message != "Knock knock!" {
		t.Error("Message not passed through")
	}
	if string(response[:n]) != "Who's there?" {
		t.Error("Response not passed through")
	}
	
	//and close
	if err := conn.Close(); err != nil {
		t.Fatal("Fatal - Can't close")
	}
	
	<-time.After(time.Second/10)
	
	if !isClosed {
		t.Error("Should be closed")
	}
}

func TestBinary(t *testing.T) {
	//setup test variables
	var isOpened,isClosed bool
	message := []byte{}
	
	//create websocket handler
	ws := New()
	ws.OnOpen(rack.Func(func(vars map[string]interface{},next func()) {
		isOpened = true
		next()
	}))
	ws.OnMessage(rack.Func(func(vars map[string]interface{},next func()) {
		message = V(vars).GetMessage().([]byte)
		V(vars).SetResponse([]byte("Who's there?"))
		next()
	}))
	ws.OnClose(rack.Func(func(vars map[string]interface{},next func()) {
		isClosed = true
		next()
	}))
	ws.ReceiveBinaryMessages()
	
	//start a web server
	server := httper.HttpConnection(":4017")
	go server.Go(ws)
	
	//connect via websocket
	conn,err := websocket.Dial("ws://localhost:4017","","http://localhost")
	if err != nil {
		t.Fatal("Fatal - Can't open")
	}
	
	if !isOpened {
		t.Error("Should be open")
	}
	
	//send a message
	if _,err := conn.Write([]byte("Knock knock!")); err != nil {
		t.Fatal("Fatal - Can't write")
	}
	
	//get a response
	response := make([]byte,512)
	n, err := conn.Read(response)
	if err != nil {
		t.Fatal("Fatal - Can't read")
	}
	
	if !bytes.Equal(message, []byte("Knock knock!")) {
		t.Error("Message not passed through")
	}
	if !bytes.Equal(response[:n], []byte("Who's there?")) {
		t.Error("Response not passed through")
	}
	
	//and close
	if err := conn.Close(); err != nil {
		t.Fatal("Fatal - Can't close")
	}
	
	<-time.After(time.Second/10)
	
	if !isClosed {
		t.Error("Should be closed")
	}
}

type TestJSONMessage struct {
	Type string `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func TestJSON(t *testing.T) {
	//setup test variables
	var isOpened,isClosed bool
	message := TestJSONMessage{}
	
	//create websocket handler
	ws := New()
	ws.OnOpen(rack.Func(func(vars map[string]interface{},next func()) {
		isOpened = true
		next()
	}))
	ws.OnMessage(rack.Func(func(vars map[string]interface{},next func()) {
		message = V(vars).GetMessage().(TestJSONMessage)
		fmt.Println(message)
		V(vars).SetResponse("Who's there?")
		next()
	}))
	ws.OnClose(rack.Func(func(vars map[string]interface{},next func()) {
		isClosed = true
		next()
	}))
	
	defaultObj := TestJSONMessage{Type:"Blank",Data:map[string]interface{}{"test":"test"}}
	ws.ReceiveJSONObjects(defaultObj)
	
	//start a web server
	server := httper.HttpConnection(":4018")
	go server.Go(ws)
	
	//connect via websocket
	conn,err := websocket.Dial("ws://localhost:4018","","http://localhost")
	if err != nil {
		t.Fatal("Fatal - Can't open")
	}
	
	if !isOpened {
		t.Error("Should be open")
	}
	
	//send a message
	sendme := TestJSONMessage{Type:"text",Data:map[string]interface{}{"message":"Knock knock!"}}
	data,err := json.Marshal(sendme)
	if err != nil {
		t.Fatal("Fatal - Can't json")
	}

	if _,err := conn.Write(data); err != nil {
		t.Fatal("Fatal - Can't write")
	}

	//get a response
	response := make([]byte,512)
	n, err := conn.Read(response)
		if err != nil {
			t.Fatal("Fatal - Can't read")
		}
	
	//response is in JSON, convert back to regular
	var response_object string
	json.Unmarshal(response[:n],&response_object)
	
	//check for message & response errors
	if message.Type != "text" {
		t.Error("Message Type Wrong - ",message.Type)
	}
	if message.Data["message"] != "Knock knock!"  {
		t.Error("Message Wrong - ",message.Data["message"])
	}
	if response_object != "Who's there?" {
		t.Error("Response not passed through - ",response_object)
	}
	
	//and close
	if err := conn.Close(); err != nil {
		t.Fatal("Fatal - Can't close")
	}
	
	<-time.After(time.Second/10)
	
	if !isClosed {
		t.Error("Should be closed")
	}
}