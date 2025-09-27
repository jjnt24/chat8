package model

type User struct {
	ID           int64  `db:"id" json:"id"`
	Username     string `db:"username" json:"username"`
	PasswordHash string `db:"password_hash" json:"-"`
	CityID       int64  `db:"city_id" json:"city_id"`
}
