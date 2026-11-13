package endpoint

type CreateRequest struct {
	ApplicationID string `json:"-"` // from path {appID}
	URL           string `json:"url"`
}
