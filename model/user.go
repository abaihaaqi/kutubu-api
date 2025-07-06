package model

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Password string `json:"password,omitempty"`

	Books []Book `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
