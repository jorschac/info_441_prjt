package users

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//gravatarBasePhotoURL is the base URL for Gravatar image requests.
//See https://id.gravatar.com/site/implement/images/ for details
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID          int64  `json:"id"`
	Email       string `json:"-"` //never JSON encoded/decoded
	PassHash    []byte `json:"-"` //never JSON encoded/decoded
	UserName    string `json:"userName"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhotoURL    string `json:"photoURL"`
	Description string `json:"description"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Description  string `json:"description"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	Description string `json:"description"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	_, err := mail.ParseAddress(nu.Email)
	if err != nil {
		return fmt.Errorf("Not a valid email adress")
	}
	if len(nu.Password) < 6 {
		return fmt.Errorf("Password too short")
	}
	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("Passwords do not match")
	}
	if len(nu.UserName) < 1 {
		return fmt.Errorf("Username has to have at least 1 character")
	}
	if strings.Contains(nu.UserName, " ") {
		return fmt.Errorf("No spaces allowed in username")
	}
	return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	err := nu.Validate()
	if err != nil {
		return nil, err
	}

	us := &User{}

	us.FirstName = nu.FirstName
	us.LastName = nu.LastName
	us.UserName = nu.UserName
	us.Email = nu.Email
	us.Description = nu.Description

	passErr := us.SetPassword(nu.Password)
	if passErr != nil {
		return nil, fmt.Errorf("Failed to hash password")
	}

	lowermail := strings.ToLower(nu.Email)
	nowhite := strings.TrimSpace(lowermail)

	hasher := md5.New()
	hasher.Write([]byte(nowhite))
	hash := hex.EncodeToString(hasher.Sum(nil))

	url := []string{"https://www.gravatar.com/avatar/", hash}
	us.PhotoURL = strings.Join(url, "")

	return us, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put between the names. If both are missing,
//this returns an empty string
func (u *User) FullName() string {
	if u.FirstName != "" && u.LastName != "" {
		concat := []string{u.FirstName, " ", u.LastName}
		return strings.Join(concat, "")
	} else if u.FirstName == "" || u.LastName == "" {
		concat := []string{u.FirstName, "", u.LastName}
		return strings.Join(concat, "")
	} else {
		return ""
	}
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return fmt.Errorf("error generating bcrypt hash: %v\n", err)
	}
	u.PassHash = hash
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	if err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password)); err != nil {
		return fmt.Errorf("password doesn't match stored hash!\n")
	} else {
		return nil
	}
}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	if updates.Description != "" {
		u.Description = updates.Description
	} else {
		return fmt.Errorf("Invalid description")
	}

	return nil
}
