package endpoint

type CreateEndpoint struct {
	ApplicationID string `json:"application_id"`
	Url           string `json:"url"`
}
