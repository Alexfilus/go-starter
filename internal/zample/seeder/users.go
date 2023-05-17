package seeder

import (
	"time"

	"github.com/otyang/go-pkg/utils"
	"github.com/otyang/yasante/internal/zample/entity"
)

var Users = []entity.User{
	{
		ID:        utils.RandomID(10),
		Email:     pS("user1@domain.com"),
		Phone:     pS("+987654321"),
		FullName:  pS("Google Man"),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        utils.RandomID(10),
		Email:     pS("user2@domain.com"),
		Phone:     pS("+987654320"),
		FullName:  pS("Google Woman"),
		IsActive:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        utils.RandomID(10),
		Email:     pS("user3@domain.com"),
		Phone:     nil,
		FullName:  nil,
		IsActive:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

// pointerString
func pS(s string) *string {
	return &s
}
