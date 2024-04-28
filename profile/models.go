package profile

import "time"

type Profile struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Email     string    `json:"email"`
	UserName  string    `json:"username"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProfileResponse struct {
	Profile *Profile
	Message string
}

var UpdateProfileReq struct {
	Bio      string `json:"bio"`
	Username string `json:"username"`
}
