package message

import "encoding/json"

type CreateMessage struct {
	ApplicationID string          `json:"-"`
	Payload       json.RawMessage `json:"payload"`
}
