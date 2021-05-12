package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var ps = PubSub{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// subscribeHandler creates new subscriptions for each new connection
func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	// We don't have to worry about this when running locally
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true

	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	id := uuid.New().String()
	ps.Subscribe(id, conn)
	fmt.Println("New subscription incoming, total: ", len(ps.Subs))
	conn.WriteMessage(1, []byte("subscriber id "+id))
}

// publishHandler publishes new messages to subscribers
func publishHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	var message Message
	err := json.Unmarshal(b, &message)
	if err != nil {
		log.Fatal(err)
	}

	_, err = ps.Publish(message.Message)
	if err != nil {
		log.Fatal(err)
	}
}
