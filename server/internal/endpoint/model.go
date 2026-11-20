package endpoint

type CreateRequest struct {
	ApplicationID string `json:"-"`
	URL           string `json:"url"`
}
