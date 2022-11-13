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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}
