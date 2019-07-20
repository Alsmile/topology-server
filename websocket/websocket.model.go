package websocket

// WsMsg websocket消息结构
type WsMsg struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}
