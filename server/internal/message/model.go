package message

import "encoding/json"

type CreateMessage struct {
	ApplicationID string          `json:"-"` // from path {appID}
	Payload       json.RawMessage `json:"payload"`
}
