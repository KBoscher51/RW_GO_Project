package models

type Victime struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Password  string `json:"password" db:"password"`
	CreatedAt string `json:"created_at" db:"created_at"`
}
