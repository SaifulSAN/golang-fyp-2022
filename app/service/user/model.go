package user

import (
	"golang.org/x/crypto/bcrypt"
)

type UserStandard struct {
	UserName                      string
	UserPhone                     string
	UserEmail                     string
	UserPassword                  string
	UserPrimaryEmergencyContact   string
	UserSecondaryEmergencyContact string
}

type UserOrganization struct {
	UserName            string
	UserPhone           string
	UserEmail           string
	UserPassword        string
	OrgSecondaryContact string
}

type UserLogin struct {
	UserEmail    string
	UserPassword string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckHashPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
