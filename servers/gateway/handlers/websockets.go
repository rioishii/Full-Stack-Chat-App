package handlers

import (
	"log"
	"net/http"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-rioishii/servers/gateway/sessions"
	"github.com/gorilla/websocket"
)

// Control messages for websocket
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

func (ctx *HandlerCtx) WebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	if !upgrader.CheckOrigin(r) {
		http.Error(w, "Websocket Connection Refused", 403)
	}
	ss := &SessionState{}
	_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, ss)
	if err != nil {
		http.Error(w, "User is not authenticated", http.StatusUnauthorized)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to open websocket connection", 401)
		return
	}

	ctx.Notifier.InsertConnection(conn, ss.User.ID)

	go (func(conn *websocket.Conn) {
		defer conn.Close()
		defer ctx.Notifier.RemoveConnection(ss.User.ID)

		for {
			messageType, p, err := conn.ReadMessage()

			if messageType == TextMessage || messageType == BinaryMessage {
				ctx.Notifier.WriteToAllConnections(TextMessage, append([]byte("Hello from server: "), p...))
			} else if messageType == CloseMessage {
				log.Println("Close message received.")
				break
			} else if err != nil {
				log.Println("Error reading message.")
				break
			}
		}

	})(conn)
}
