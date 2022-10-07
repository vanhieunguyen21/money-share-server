package request

import "money_share/pkg/dto"

type GroupCreationRequest struct {
	Group     dto.GroupDTO `json:"group"`
	CreatorID uint         `json:"creatorID"`
}
