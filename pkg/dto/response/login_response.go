package response

import "money_share/pkg/dto"

type LoginResponse struct {
	dto.UserDTO
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
