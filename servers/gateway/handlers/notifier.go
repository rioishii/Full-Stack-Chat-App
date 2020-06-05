package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

type Notifier struct {
	Connections map[int64]*websocket.Conn
	lock        sync.Mutex
}

func NewNotifier() *Notifier {
	return &Notifier{}
}

type Message map[string]interface{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		if r.Header.Get("Origin") != "https://rioishii.me" {
			return false
		}
		return true
	},
}

func (n *Notifier) InsertConnection(conn *websocket.Conn, userID int64) {
	n.lock.Lock()
	defer n.lock.Unlock()
	if n.Connections == nil {
		n.Connections = make(map[int64]*websocket.Conn)
	}
	n.Connections[userID] = conn
}

func (n *Notifier) RemoveConnection(userID int64) {
	n.lock.Lock()
	defer n.lock.Unlock()
	// delete socket connection
	delete(n.Connections, userID)
}

func (n *Notifier) WriteToAllConnections(messageType int, data []byte) error {
	var writeError error
	for id, conn := range n.Connections {
		writeError = conn.WriteMessage(messageType, data)
		if writeError != nil {
			n.RemoveConnection(id)
			conn.Close()
			return writeError
		}
	}
	return nil
}

func (n *Notifier) WriteToConnection(messageType int, data []byte, userIDs []int64) error {
	var writeError error
	for _, id := range userIDs {
		if conn, ok := n.Connections[id]; ok {
			writeError = conn.WriteMessage(messageType, data)
			if writeError != nil {
				n.RemoveConnection(id)
				conn.Close()
				return writeError
			}
		}
	}
	return nil
}

// go routine executed in main for sending message to the proper websockets
func (n *Notifier) NotifyWebSockets(msgs <-chan amqp.Delivery) {
	for m := range msgs {
		n.lock.Lock()
		log.Printf("Received a message: %s", m.Body)

		var message Message
		buf := bytes.NewBuffer(m.Body)
		decoder := json.NewDecoder(buf)
		err := decoder.Decode(&message)
		if err != nil {
			log.Printf("Error receiving from queue: %s", err.Error())
		}
		if message["userIDs"] == nil {
			err = n.WriteToAllConnections(TextMessage, m.Body)
			if err != nil {
				log.Printf("Error writing from queue: %s", err.Error())
			}
		} else {
			idArr := message["userIDs"].([]interface{})
			convertedIdArr := make([]float64, len(idArr))
			convertedIdArr2 := make([]int64, len(idArr))
			for i := range idArr {
				convertedIdArr[i] = idArr[i].(float64)
			}
			for i := range idArr {
				convertedIdArr2[i] = int64(convertedIdArr[i])
			}
			err = n.WriteToConnection(TextMessage, m.Body, convertedIdArr2)
			if err != nil {
				log.Printf("Error writing from queue: %s", err.Error())
			}
		}
		n.lock.Unlock()
	}
}
