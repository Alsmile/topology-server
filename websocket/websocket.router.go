package websocket

import (
	"encoding/json"
	"topology/middlewares"

	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	"github.com/rs/zerolog/log"
)

// Route Setup websocket.
func Route(app *iris.Application) {
	// create our echo websocket server
	ws := websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})
	ws.OnConnection(wsConnection)

	// register the server on an endpoint.
	// see the inline javascript code in the websockets.html,
	// this endpoint is used to connect to the server.
	app.Get("/ws", ws.Handler())
}

func wsConnection(conn websocket.Connection) {
	conn.OnMessage(func(data []byte) {
		var msg WsMsg
		err := json.Unmarshal(data, &msg)
		if err != nil {
			log.Error().Caller().Err(err).Msgf("Error on wsConnection: %s", string(data))
			return
		}
		switch msg.Event {
		case "token":
			middlewares.ParseJwt(conn.Context(), msg.Data)
		}

		log.Debug().Caller().Msgf("Websocket OnMessage: userId=%s, data=%s", conn.Context().Values().GetString("uid"), string(data))
	})

	info := make(map[string]string)
	info["remoteAddr"] = conn.Context().RemoteAddr()
	jsonStr, err := json.Marshal(info)
	if err == nil {
		conn.EmitMessage([]byte(jsonStr))
	}
}
