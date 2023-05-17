package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            string    `json:"id" bun:",pk"`
	Email         *string   `json:"email" bun:",unique"`
	Phone         *string   `json:"phoneNumber" bun:",unique"`
	FullName      *string   `json:"fullName" `
	IsActive      bool      `json:"isActive"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
