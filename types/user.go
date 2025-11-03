package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

const (
	COST            = 12
	minFirstNameLen = 3
	minLastNameLen  = 3
	minPasswordLen  = 8
)

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type User struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // these are called struct tags, it will omit the ID from json if found empty
	FirstName   string        `bson:"firstName" json:"firstName"`
	LastName    string        `bson:"lastName" json:"lastName"`
	Email       string        `bson:"email" json:"email"`
	EncPassword string        `bson:"EncPassword" json:"-"`
}

func (p UpdateUserParams) ToBSON() bson.M {
	values := bson.M{}
	if len(p.FirstName) > 0 {
		values["firstName"] = p.FirstName
	}
	if len(p.LastName) > 0 {
		values["lastName"] = p.LastName
	}
	return values
}
func IsValidPassword(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil

}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func (params CreateUserParams) Validate() map[string]string {
	errs := map[string]string{}
	if len(params.FirstName) < minFirstNameLen {
		errs["firstName"] = (fmt.Sprintf("first name length should be at least %d chars", minFirstNameLen))
	}
	if len(params.LastName) < minLastNameLen {
		errs["lastName"] = (fmt.Sprintf("last name length should be at least %d chars", minLastNameLen))
	}
	if len(params.Password) < minPasswordLen {
		errs["password"] = (fmt.Sprintf("password length should be at least %d chars", minPasswordLen))
	}
	if !isEmailValid(params.Email) {
		errs["email"] = ("not a valid email")
	}
	if len(errs) > 0 {
		return errs
	}
	return errs
}

func NewUserFromParams(params CreateUserParams) (*User, error) {

	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), 12)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:   params.FirstName,
		LastName:    params.LastName,
		Email:       params.Email,
		EncPassword: string(encpw),
	}, nil
}
