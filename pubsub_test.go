package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestUnsubscribe(t *testing.T) {
	cases := []struct {
		ps       *PubSub
		id       string
		expected func() ([]Subscription, error)
	}{
		{
			ps: &PubSub{
				Subs: []Subscription{
					{
						ID: "2345",
					},
					{
						ID: "6789",
					},
				},
			},
			id: "2345",
			expected: func() ([]Subscription, error) {
				return []Subscription{{ID: "6789"}}, nil
			},
		},
		{
			ps: &PubSub{
				Subs: []Subscription{
					{
						ID: "6789",
					},
				},
			},
			id: "2345",
			expected: func() ([]Subscription, error) {
				errText := fmt.Sprintf("could not unsubscribe id %v not found", "2345")
				return []Subscription{{ID: "6789"}}, errors.New(errText)
			},
		},
		{
			ps: &PubSub{
				Subs: []Subscription{},
			},
			id: "2345",
			expected: func() ([]Subscription, error) {
				errText := fmt.Sprintf("could not unsubscribe id %v not found", "2345")
				return []Subscription{}, errors.New(errText)
			},
		},
	}

	for _, c := range cases {
		actualErr := c.ps.Unsubscribe(c.id)
		subs, expectedErr := c.expected()

		assert.Equal(t, actualErr, expectedErr)
		assert.Equal(t, c.ps.Subs, subs)
	}
}

func TestSubscribe(t *testing.T) {
	cases := []struct {
		sub      *Subscription
		ps       *PubSub
		expected []Subscription
	}{
		{
			sub: &Subscription{
				ID: "1234",
			},
			ps: &PubSub{
				Subs: []Subscription{
					{
						ID: "5678",
					},
				},
			},
			expected: []Subscription{
				{
					ID: "5678",
				},
				{
					ID: "1234",
				},
			},
		},
		{
			sub: &Subscription{
				ID: "1234",
			},
			ps: &PubSub{
				Subs: []Subscription{},
			},
			expected: []Subscription{
				{
					ID: "1234",
				},
			},
		},
	}

	for _, c := range cases {
		c.ps.Subscribe(c.sub.ID, nil)

		assert.Equal(t, c.ps.Subs, c.expected)
	}
}

// This test just checks that the correct number of subs are published to
func TestPublish(t *testing.T) {
	ps := &PubSub{}

	s := httptest.NewServer(http.HandlerFunc(subscribeHandler))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}

	newSub := Subscription{
		ID:         "12345",
		Connection: ws,
	}

	ps.Subs = append(ps.Subs, newSub)

	count, err := ps.Publish([]byte(`{"message":"hey"`))
	if err != nil {
		t.Fatalf("%v", err)
	}

	assert.Equal(t, count, 1)

	secondSub := Subscription{
		ID:         "67890",
		Connection: ws,
	}

	ps.Subs = append(ps.Subs, secondSub)

	count, err = ps.Publish([]byte(`{"message":"hey"`))
	if err != nil {
		t.Fatalf("%v", err)
	}

	assert.Equal(t, count, 2)
}
