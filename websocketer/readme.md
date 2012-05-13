#	Websocketer
Websocketer provides a simple rack based interface for starting up a websocket, and allows you to specify a rack based interface for dealing with websocket based events

## 	Dependencies
This uses google's Websocket code

## 	Installation
`go get github.com/HairyMezican/Middleware/websocketer`

## 	Example

    package main

	import (
		"github.com/HairyMezican/Middleware/logger"
		"github.com/HairyMezican/Middleware/websocketer"
		"github.com/HairyMezican/TheRack/rack"
	)

	var OpenWare rack.Func = func(vars rack.Vars, next func()) {
		lg := logger.Get(vars)
		lg.Print("Opened")
		next()
	}

	var OpeningMessage rack.Func = func(vars rack.Vars, next func()) {
		websocketer.SendBasicMessage(vars, "Welcome!")
		next()
	}

	var MessageWare rack.Func = func(vars rack.Vars, next func()) {
		lg := logger.Get(vars)
		s, ok := websocketer.Message(vars).(*string)
		if ok {
			lg.Print("Succesful Message - \"" + *s + "\"")
		}
		next()
	}

	var CloseWare rack.Func = func(vars rack.Vars, next func()) {
		lg := logger.Get(vars)
		lg.Print("Closed")
		next()
	}

	var HttpWare rack.Func = func(vars rack.Vars, next func()) {
		rack.SetMessageString(vars, "<html><head><title>Websockets</title><script type='text/javascript'>c=new WebSocket('ws://localhost:3000');c.onmessage=function(e){c.send('Message Received')}</script></head><body>Websockets!</body></html>")
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

		conn := rack.HttpConnection(":3000")
		conn.Go(rackup)
	}
	
Opening localhost:3000 should just display "Websockets!", looking at the standard output, though, it should show you that websocket was opened, and that you received a message, and upon closing the browser window, it should tell you that the websocket was closed