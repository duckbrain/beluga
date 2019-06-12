package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            uuid.UUID    `json:"id" db:"id"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" db:"updated_at"`
	Username      nulls.String `json:"username" db:"username"` // Can ommit username and password, then is used for deploys only
	PasswordHash  string       `json:"password_hash" db:"password_hash"`
	IsAdmin       bool         `json:"is_admin" db:"is_admin"`             // Is allowed to create/edit users
	Key           string       `json:"key" db:"key"`                       // Can only be used for deploy and teardown
	DomainPattern string       `json:"domain_pattern" db:"domain_pattern"` // Regex to match allowed stack names
}

func checkPwd(pwd, hash string) (needsUpdate bool, err error) {
	if len(hash) == 0 {
		return false, errors.New("No hash set")
	}
	return false, bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
}

func genPwd(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash), err
}

func (u User) VerifyKey(key string) (needsUpdate bool, err error) {
	if key != u.Key {
		err = errors.New("Keys don't match")
	}
	return
}

func (u User) VerifyPassword(password string) (needsUpdate bool, err error) {
	return checkPwd(password, u.PasswordHash)
}

func (u *User) SetPassword(pwd string) (err error) {
	u.PasswordHash, err = genPwd(pwd)
	return
}

func (u *User) GenerateKey() (key string, err error) {
	key, err = GenerateKey(KeyLength)
	if err == nil {
		u.Key = key
	}
	return
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	v := []validate.Validator{}
	if u.Username.Valid {
		v = append(v, &validators.StringIsPresent{Field: u.Username.String, Name: "Username"})
	}

	return validate.Validate(v...), nil
}
