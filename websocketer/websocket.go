package websocketer

import (
	"code.google.com/p/go.net/websocket"
	"github.com/HairyMezican/Middleware/logger"
	"github.com/HairyMezican/TheRack/httper"
	"github.com/HairyMezican/TheRack/rack"
	"io"
	"strings"
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

func (this Middleware) Run(vars map[string]interface{}, next func()) {
	r := (httper.V)(vars).GetRequest()
	if strings.ToLower(r.Header.Get("Upgrade")) != "websocket" {
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
					if err != io.EOF {
						lg := (logger.V)(vars).Get()
						if lg != nil {
							lg.Println(err)
						}
						continue
					} else {
						break
					}
				}

				//respond to the message in a goroutine
				go func() {
					vars[messageIndex] = message
					this.onMessage.Run(vars, func() {})

					//If we have a response, send it back
					response := vars[responseIndex]
					delete(vars, responseIndex)

					if response != nil {
						this.messageType.Send(ws, response)
					}
				}()
			}
		}
		w := (httper.V)(vars).BlankResponse()
		handler.ServeHTTP(w, r)
		w.Save()
	}
}

func (this *Middleware) UseJSON() {
	this.messageType = websocket.JSON
}

func (this *Middleware) UseMessage() {
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

func (this *Middleware) OnStorage(f func() interface{}) {
	this.onStorage = f
}

type V map[string]interface{}

func (vars V) GetMessage() interface{} {
	return vars[messageIndex]
}

func (vars V) SetResponse(response interface{}) {
	vars[responseIndex] = response
}

func (vars V) SendBasicMessage(message interface{}) {
	vars.sendmessage(message, websocket.Message)
}

func (vars V) SendJSONMessage(message interface{}) {
	vars.sendmessage(message, websocket.JSON)
}

func (vars V) sendmessage(message interface{}, c websocket.Codec) {
	ws := vars.GetSocket()
	c.Send(ws, message)
}

func (vars V) GetSocket() *websocket.Conn {
	ws, ok := vars[websocketIndex].(*websocket.Conn)
	if !ok {
		panic("Can't find websocket")
	}
	return ws
}
