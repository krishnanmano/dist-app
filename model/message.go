package model

import (
	"sync"
	"time"
)

var mut sync.Mutex
var message = make([]Message, 0)

type Message struct {
	Msg       string `json:"message"`
	CreatedAt int64  `json:"created_at,omitempty"`
}

func (msg Message) SaveMessage() {
	mut.Lock()
	defer mut.Unlock()
	msg.CreatedAt = time.Now().UnixMilli()
	message = append(message, msg)
}

func (msg Message) GetMessages() []Message {
	return message

}
