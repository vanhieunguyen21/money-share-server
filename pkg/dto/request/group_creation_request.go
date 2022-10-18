package request

type GroupCreationRequest struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}
