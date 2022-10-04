package response

type CheckUsernameResponse struct {
	Username    string `json:"username"`
	Requirement bool   `json:"requirement"`
	Available   bool   `json:"available"`
	Message     string `json:"message"`
}
