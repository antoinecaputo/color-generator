package main

import (
	"endiatx/lib/golog"
	"endiatx/lib/gosocket"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func main() {
	logWebsocketConnection := &gosocket.WebSocketConnection{
		Connection: nil,
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}

	go gosocket.ServeSocket(logWebsocketConnection, 8080)
	defer logWebsocketConnection.Close()

	logger, err := golog.NewLogger(golog.LogLvlDebug)
	if err != nil {
		log.Fatalln(err)
	}

	err = logger.FctAddOutput(logWebsocketConnection)
	if err != nil {
		log.Fatalln(err)
	}

	logger.Start()
}
