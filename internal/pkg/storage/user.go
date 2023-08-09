package storage

import (
	"context"
	"crypto/rsa"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/zltl/xoidc/internal/pkg/db"
	"golang.org/x/text/language"
)

type User struct {
	ID                string
	Username          string
	Password          string
	FirstName         string
	LastName          string
	Email             string
	EmailVerified     bool
	Phone             string
	PhoneVerified     bool
	PreferredLanguage language.Tag
	IsAdmin           bool
}

type Service struct {
	keys map[string]*rsa.PublicKey
}

type UserStore interface {
	GetUserByID(string) *User
	GetUserByUsername(string) *User
	ExampleClientID() string
}

type userStore struct {
	users map[string]*User
	DB    *db.Store
}

func NewUserStore(issuer string) UserStore {
	hostname := strings.Split(strings.Split(issuer, "://")[1], ":")[0]
	return &userStore{
		users: map[string]*User{
			"id1": {
				ID:                "id1",
				Username:          "test-user@" + hostname,
				Password:          "verysecure",
				FirstName:         "Test",
				LastName:          "User",
				Email:             "test-user@zitadel.ch",
				EmailVerified:     true,
				Phone:             "",
				PhoneVerified:     false,
				PreferredLanguage: language.German,
				IsAdmin:           true,
			},
			"id2": {
				ID:                "id2",
				Username:          "test-user2",
				Password:          "verysecure",
				FirstName:         "Test",
				LastName:          "User2",
				Email:             "test-user2@zitadel.ch",
				EmailVerified:     true,
				Phone:             "",
				PhoneVerified:     false,
				PreferredLanguage: language.German,
				IsAdmin:           false,
			},
		},
	}
}

// ExampleClientID is only used in the example server
func (u userStore) ExampleClientID() string {
	return "service"
}

func (u userStore) GetUserByID(id string) *User {
	us, err := u.DB.GetUserByID(context.TODO(), id)
	if err != nil {
		return nil
	}
	sid := snowflake.ID(us.ID)
	return &User{
		ID:            sid.Base64(),
		Username:      us.Username,
		Password:      us.Password,
		FirstName:     us.GivenName,
		LastName:      us.FamilyName,
		Email:         us.Email,
		EmailVerified: us.EmailVerified,
		Phone:         us.PhoneNumber,
		PhoneVerified: us.PhoneNumberVerified,
	}
}

func (u userStore) GetUserByUsername(username string) *User {
	us, err := u.DB.GetUserByUsername(context.TODO(), username)
	if err != nil {
		return nil
	}
	sid := snowflake.ID(us.ID)
	return &User{
		ID:            sid.Base64(),
		Username:      us.Username,
		Password:      us.Password,
		FirstName:     us.GivenName,
		LastName:      us.FamilyName,
		Email:         us.Email,
		EmailVerified: us.EmailVerified,
		Phone:         us.PhoneNumber,
		PhoneVerified: us.PhoneNumberVerified,
	}
}
