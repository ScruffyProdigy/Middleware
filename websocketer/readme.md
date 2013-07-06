#	Websocketer
Websocketer provides a simple rack based interface for starting up a websocket, and allows you to specify a rack based interface for dealing with websocket based events

## 	Installation
`go get github.com/ScruffyProdigy/Middleware/websocketer`

##  Documentation
http://godoc.org/github.com/ScruffyProdigy/Middleware/websocketer

## Usage

* Call websocketer.New() to get a Controller object
* Decide what format you want to receive your messages in
	* By default, you will accept text messages (as `string`)
		* You can call the controller's ReceiveTextMessages() to make it explicit
	* If you want to accept binary data (as `[]byte`), call the controller's ReceiveBinaryMessages()
	* If you want to accept JSON data, you will need to call the controller's ReceiveJSONObjects()
		* You will need to send it an object so it can know what format to give you your data in
* Handle messages being sent from the user by calling the controller's OnMessage() and sending it a Middleware that will handle the message
	* You can call GetMessage() as a vars function to get the message that was sent to you
		* It will give you the message as an interface{}, but the underlying data type will be determined by whatever you specified in the previous step
	* You can call SetResponse() as a vars function to get specify the response that should be sent back
* To handle users coming and going, you can call the controller's OnOpen() and OnClose() to specify Middleware to handle the setup and cleanup of user data
	* It is a good idea to store the vars somewhere accessible to be able to use the next step
* To send a user a message when there is no callback going on, use the vars functions SendJSONMessage() and SendBasicMessage()
* I think that covers all of the basic functionality, but in case I missed something, GetSocket() has been provided as a safety valve
	* This will give you access to the user's websocket connection directly, see http://godoc.org/code.google.com/p/go.net/websocket for more information

## 	Example

    package main

	import (
		"github.com/ScruffyProdigy/Middleware/logger"
		"github.com/ScruffyProdigy/Middleware/websocketer"
		"github.com/ScruffyProdigy/TheRack/httper"
		"github.com/ScruffyProdigy/TheRack/rack"
	)

	var OpenWare rack.Func = func(vars map[string]interface{}, next func()) {
		lg := (logger.V)(vars).Get()
		lg.Print("Opened")
		next()
	}

	var OpeningMessage rack.Func = func(vars map[string]interface{}, next func()) {
		(websocketer.V)(vars).SendBasicMessage("Welcome!")
		next()
	}

	var MessageWare rack.Func = func(vars map[string]interface{}, next func()) {
		lg := (logger.V)(vars).Get()
		s, ok := (websocketer.V)(vars).GetMessage().(*string)
		if ok {
			lg.Print("Succesful Message - \"" + *s + "\"")
		}
		next()
	}

	var CloseWare rack.Func = func(vars map[string]interface{}, next func()) {
		lg := (logger.V)(vars).Get()
		lg.Print("Closed")
		next()
	}

	var HttpWare rack.Func = func(vars map[string]interface{}, next func()) {
		(httper.V)(vars).SetMessageString("<html><head><title>Websockets</title><script type='text/javascript'>c=new WebSocket('ws://localhost:3000');c.onmessage=function(e){c.send('Message Received')}</script></head><body>Websockets!</body></html>")
	}

	func main() {
		ws := websocketer.New()
		ws.OnOpen(OpenWare)
		ws.OnOpen(OpeningMessage)
		ws.OnMessage(MessageWare)
		ws.OnClose(CloseWare)

		rackup := rack.New()
		rackup.Add(logger.StandardLogger)
		rackup.Add(ws)
		rackup.Add(HttpWare)

		conn := httper.HttpConnection(":3000")
		conn.Go(rackup)
	}
	
	
Opening localhost:3000 should just display "Websockets!", looking at the standard output, though, it should show you that websocket was opened, and that you received a message, and upon closing the browser window, it should tell you that the websocket was closed