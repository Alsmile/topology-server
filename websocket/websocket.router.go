package websocket

import (
	"encoding/json"
	"topology/middlewares"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"github.com/rs/zerolog/log"
)

// Route Setup websocket.
func Route(app *iris.Application) {
	// create our echo websocket server
	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: OnMessage,
	})

	// register the server on an endpoint.
	// see the inline javascript code in the websockets.html,
	// this endpoint is used to connect to the server.
	app.Get("/ws", websocket.Handler(ws))
}

// OnMessage websocket 消息处理函数
func OnMessage(nsConn *websocket.NSConn, msg websocket.Message) error {
	ctx := websocket.GetContext(nsConn.Conn)

	var wsMsg WsMsg
	err := json.Unmarshal(msg.Body, &wsMsg)
	if err != nil {
		log.Error().Caller().Err(err).Msgf("Error on wsConnection: %s", msg.Body)
		return nil
	}
	switch msg.Event {
	case "token":
		middlewares.ParseJwt(ctx, wsMsg.Data)
	}

	log.Debug().Caller().Msgf("Websocket OnMessage: userId=%s, data=%s", ctx.Values().GetString("uid"), msg.Body)

	return nil
}
