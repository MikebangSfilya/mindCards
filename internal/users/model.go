package users

type Users struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
	UserId   int    `json:"user_id" db:"user_id"`
}
