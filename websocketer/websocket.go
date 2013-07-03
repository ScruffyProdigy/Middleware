/*
	websocketer provides a middleware that will easily add websocet technology into your web app
*/
package websocketer

import (
	"code.google.com/p/go.net/websocket"
	"github.com/ScruffyProdigy/Middleware/logger"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"io"
	"reflect"
	"strings"
)

const (
	messageIndex   = "Message"
	responseIndex  = "Response"
	websocketIndex = "Websocket"
)

// New() returns a blank Websocket Controller
func New() *Controller {
	this := new(Controller)
	this.onOpen = rack.New()
	this.onMessage = rack.New()
	this.onClose = rack.New()
	this.ReceiveTextMessages()
	return this
}

// Controller() controls websocket flow
type Controller struct {
	messageType                websocket.Codec
	onStorage                  func() interface{} //should return a pointer to whatever you want the messages stored in
	onOpen, onMessage, onClose *rack.Rack
}

// Run() implements the rack.Middleware interface
func (this Controller) Run(vars map[string]interface{}, next func()) {
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
						(logger.V)(vars).Println(err)
						continue
					} else {
						break
					}
				}

				//respond to the message in a goroutine
				go func() {
					vars[messageIndex] = reflect.Indirect(reflect.ValueOf(message)).Interface()
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

// ReceiveTextMessages() is the default, but calling it will make explicit that we are expecting to receive strings from the client
func (this *Controller) ReceiveTextMessages() {
	this.messageType = websocket.Message
	this.onStorage = func() interface{} {
		var m string
		return &m
	}
}

// ReceiveBinaryMessages() declares that we are expecting to receive a slice of bytes from the client
func (this *Controller) ReceiveBinaryMessages() {
	this.messageType = websocket.Message
	this.onStorage = func() interface{} {
		var m []byte
		return &m
	}
}

// ReceiveJSONObjects() declares that we are expecting to receive JSON data from the client,
// and want it to be stored in an object similar to the example we are sending
func (this *Controller) ReceiveJSONObjects(example interface{}) {
	this.messageType = websocket.JSON
	t := reflect.TypeOf(example)

	this.onStorage = func() interface{} {
		return reflect.New(t).Interface()
	}
}

// OnOpen() gives you an opportunity to handle setup when a client connects
func (this Controller) OnOpen(m rack.Middleware) {
	this.onOpen.Add(m)
}

// OnMessage() gives you an opportunity to react when the client sends you data
func (this Controller) OnMessage(m rack.Middleware) {
	this.onMessage.Add(m)
}

// OnClose() gives you an opportunity to clean up when a client disconnects
func (this Controller) OnClose(m rack.Middleware) {
	this.onClose.Add(m)
}

// V is a type you can cast your vars to in order to access the following functions
type V map[string]interface{}

//During OnMessage() you can call GetMessage() to get the data that the user sent.
//It is sent as an interface{} but the underlying type will depend on which of the ReceiveXxxxx() functions you used during setup
func (vars V) GetMessage() interface{} {
	return vars[messageIndex]
}

//During OnMessage() you can call SetResponse to set what kind of message you are sending back to the client
func (vars V) SetResponse(response interface{}) {
	vars[responseIndex] = response
}

//SendBasicMessage() can be used at any time to send a string or byte slice back to the client
func (vars V) SendBasicMessage(message interface{}) {
	vars.sendmessage(message, websocket.Message)
}

//SendJSONMessage() can be used at any time to send a JSON object back to the client
func (vars V) SendJSONMessage(message interface{}) {
	vars.sendmessage(message, websocket.JSON)
}

func (vars V) sendmessage(message interface{}, c websocket.Codec) {
	ws := vars.GetSocket()
	c.Send(ws, message)
}

// You shouldn't need this, but it has been included as a safety valve:
// GetSocket() will give you direct access to the underlying websocket connection
func (vars V) GetSocket() *websocket.Conn {
	ws, ok := vars[websocketIndex].(*websocket.Conn)
	if !ok {
		panic("Can't find websocket")
	}
	return ws
}
