package websocketer

import (
	websocket "code.google.com/p/go.net/websocket"
	"github.com/HairyMezican/TheRack/rack"
)

const (
	messageIndex   = "Message"
	responseIndex  = "Response"
	websocketIndex = "Websocket"
)

func New() *Middleware {
	this := new(Middleware)
	this.messageType = websocket.Message
	this.onStorage = func() interface{} {
		var m string
		return &m
	}
	this.onOpen = rack.New()
	this.onMessage = rack.New()
	this.onClose = rack.New()
	return this
}

type Middleware struct {
	messageType                websocket.Codec
	onStorage                  func() interface{} //should return a pointer to whatever you want the messages stored in
	onOpen, onMessage, onClose *rack.Rack
}

func (this Middleware) Run(vars rack.Vars, next func()) {
	r := rack.GetRequest(vars)
	if r.Header.Get("Upgrade") != "WebSocket" {
		//if it wasn't a websocket request, ignore it
		next()
	} else {
		var handler websocket.Handler = func(ws *websocket.Conn) {
			vars[websocketIndex] = ws
			this.onOpen.Run(vars, func() {})
			defer func() {
				this.onClose.Run(vars, func() {})
				ws.Close()
			}()
			for {
				//Message Loop

				//Get the message from the client
				message := this.onStorage()
				err := this.messageType.Receive(ws, message)

				//If there are no messages, we're done here
				if err != nil {
					break
				}

				//respond to the message in a goroutine
				go func() {
					vars[messageIndex] = message
					this.onMessage.Run(vars, func() {})

					//If we have a response, send it back
					response := vars.Clear(responseIndex)
					if response != nil {
						this.messageType.Send(ws, response)
					}
				}()
			}
		}
		w := rack.BlankResponse(vars)
		handler.ServeHTTP(w, r)
		w.Save()
	}
}

func (this Middleware) UseJSON() {
	this.messageType = websocket.JSON
}

func (this Middleware) UseMessage() {
	this.messageType = websocket.Message
}

func (this Middleware) OnOpen(m rack.Middleware) {
	this.onOpen.Add(m)
}

func (this Middleware) OnMessage(m rack.Middleware) {
	this.onMessage.Add(m)
}

func (this Middleware) OnClose(m rack.Middleware) {
	this.onClose.Add(m)
}

func (this Middleware) OnStorage(f func() interface{}) {
	this.onStorage = f
}

func Message(vars rack.Vars) interface{} {
	return vars[messageIndex]
}

func SetResponse(vars rack.Vars, response interface{}) {
	vars[responseIndex] = response
}

func SendBasicMessage(vars rack.Vars, message interface{}) {
	sendmessage(vars, message, websocket.Message)
}

func SendJSONMessage(vars rack.Vars, message interface{}) {
	sendmessage(vars, message, websocket.JSON)
}

func sendmessage(vars rack.Vars, message interface{}, c websocket.Codec) {
	ws, ok := vars[websocketIndex].(*websocket.Conn)
	if !ok {
		panic("Can't find websocket")
	}

	c.Send(ws, message)
}