package common

type TokenPayload struct {
	UserId   int    `json:"user_id"`
	Email    string `json:"email"`
	SubToken string `json:"sub_token"`
}

func (t TokenPayload) GetUserId() int {
	return t.UserId
}

func (t TokenPayload) GetSubToken() string {
	return t.SubToken
}

func (t TokenPayload) GetEmail() string {
	return t.Email
}
