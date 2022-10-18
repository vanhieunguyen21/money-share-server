package request

type UpdateUserRequest struct {
	DisplayName     *string `json:"displayName"`
	Password        *string `json:"password"`
	PhoneNumber     *string `json:"phoneNumber"`
	EmailAddress    *string `json:"emailAddress"`
	DateOfBirth     *string `json:"dateOfBirth"`
}
