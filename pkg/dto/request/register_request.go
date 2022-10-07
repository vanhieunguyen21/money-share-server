package request

import "money_share/pkg/dto"

type RegisterRequest struct {
	dto.UserDTO
	Password    string `json:"password"`
}
