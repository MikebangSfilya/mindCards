package users

type Users struct {
	Email    string `json:"email"`
	Password string `json:"-"`
	UserId   int    `json:"user_id"`
}
