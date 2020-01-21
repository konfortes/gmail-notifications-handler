package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/pubsub"
)

var (
	topic *pubsub.Topic

	// Messages received by this instance.
	messagesMu sync.Mutex
	messages   []string

	// token is used to verify push requests. TODO: panic if not present
	token = os.Getenv("PUBSUB_VERIFICATION_TOKEN")
)

const maxMessages = 10

type pushRequest struct {
	Message struct {
		Attributes map[string]string
		Data       []byte
		ID         string `json:"message_id"`
	}
	Subscription string
}

func pushHandler(w http.ResponseWriter, req *http.Request) {
	// Verify the token.
	if req.URL.Query().Get("token") != token {
		http.Error(w, "Bad token", http.StatusBadRequest)
	}
	msg := &pushRequest{}
	if err := json.NewDecoder(req.Body).Decode(msg); err != nil {
		http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
		return
	}

	messagesMu.Lock()
	defer messagesMu.Unlock()
	// Limit to ten.
	messages = append(messages, string(msg.Message.Data))
	if len(messages) > maxMessages {
		messages = messages[len(messages)-maxMessages:]
	}
}

// test deploy
func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
