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
		"github.com/HairyMezican/TheRack/httper"
		"github.com/HairyMezican/TheRack/rack"
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