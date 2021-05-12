package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

// PubSub maintains a list of current subscriptions
type PubSub struct {
	Subs []Subscription
}

type Message struct {
	Message json.RawMessage `json:"message"`
}

// Subscription has both a unique ID and a websocket connection
type Subscription struct {
	ID         string
	Connection *websocket.Conn
}

// Subscribe creates a new subscription and adds it to the current list
func (ps *PubSub) Subscribe(id string, conn *websocket.Conn) {
	newSub := Subscription{
		Connection: conn,
		ID:         id,
	}

	ps.Subs = append(ps.Subs, newSub)
}

// Unsubscribe searches for a subscription by id and removes it from the current list
func (ps *PubSub) Unsubscribe(id string) error {
	found := false
	for dex, sub := range ps.Subs {
		if sub.ID == id {
			found = true
			ps.Subs = append(ps.Subs[:dex], ps.Subs[dex+1:]...)
		}
	}

	if !found {
		errText := fmt.Sprintf("could not unsubscribe id %v not found", id)
		return errors.New(errText)
	}

	return nil
}

// Publish iterates through the current subscriptions, and writes a message to each one
func (ps *PubSub) Publish(message []byte) (int, error) {
	subscriptions := ps.Subs
	count := 0

	for _, sub := range subscriptions {
		fmt.Printf("Sending to client id %s message is %s \n", sub.ID, string(message))
		err := sub.Send(message)
		if err != nil {
			return count, err
		}

		count++
	}

	return count, nil
}

func (sub *Subscription) Send(data []byte) error {
	err := sub.Connection.WriteMessage(1, data)
	if err != nil {
		return err
	}

	return nil
}
