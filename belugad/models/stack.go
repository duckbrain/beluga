package models

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type Stack struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Name      string    `json:"name" db:"name"`
	Data      string    `json:"data" db:"data"`
}

const KeyLength = 70

func GenerateKey(length int) (string, error) {
	b := make(([]byte), length)
	n, err := rand.Read(b)
	if n < length {
		return "", errors.New("Not enough random bytes")
	}
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(b), nil
}

// String is not required by pop and may be deleted
func (s Stack) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Stacks is not required by pop and may be deleted
type Stacks []Stack

// String is not required by pop and may be deleted
func (s Stacks) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Stack) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.Name, Name: "Name"},
		validate.ValidatorFunc(func(errors *validate.Errors) {
			if strings.HasPrefix(s.Name, "beluga") {
				errors.Add("Name", "Cannot have \"beluga\" prefix.")
			}
		}),
	), nil
}
