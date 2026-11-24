package endpoint

import "github.com/findardi/forflow-webhooks/server/db/store/sqlc"

type CreateRequest struct {
	ApplicationID string `json:"-"`
	URL           string `json:"url"`
}

type CreateResponse struct {
	Secret string                 `json:"secret"`
	Result sqlc.CreateEndpointRow `json:"result"`
}
