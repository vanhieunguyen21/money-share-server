package response

import "money_share/pkg/dto"

type LoginResponse struct {
	dto.UserDTO
	Token string `json:"token"`
}
