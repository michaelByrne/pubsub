package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestPublishHandler(t *testing.T) {
	jsonStr := []byte(`{"message":"hello"}`)
	req, _ := http.NewRequest("POST", "/publish", bytes.NewBuffer(jsonStr))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(publishHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSubscribeHandler(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(subscribeHandler))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	_, m, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	assert.Contains(t, string(m), "subscriber id")
}
