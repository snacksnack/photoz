package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"../hash"
	"../rand"
)

var (
	// ErrNotFound returned when resource cannot be found in the database
	ErrNotFound        = errors.New("models: resource not found.")
	ErrInvalidID       = errors.New("models: ID must be > 0")
	ErrInvalidPassword = errors.New("models: incorrect password provided.")
)

const userPwPepper = "wtul91.5"
const hmacSecretKey = "progressivestereo"

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `sql:"-"` //telling gorm to ignore this - don't store in db
	PasswordHash string `gorm:"not null"`
	Remember     string `sql:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// look up user by id.
type UserService struct {
	db   *gorm.DB
	hmac hash.HMAC
}

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	return &UserService{
		db:   db,
		hmac: hmac,
	}, nil
}

// create user
func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(pwBytes), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	//not entirely necessary, but could help preventing raw password from appearing in logs
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		} else {
			user.Remember = token
		}
		user.RememberHash = us.hmac.Hash(token)
	}
	return us.db.Create(user).Error
}

// update user (expects complete User object)
func (us *UserService) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}
	return us.db.Save(user).Error
}

// delete user with provided ID
func (us *UserService) Delete(id uint) error {
	//gorm deletes all records in table if id == 0
	if id == 0 {
		return ErrInvalidID
	}

	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

// look up user by id
func (us *UserService) ById(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// look up user by email
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// looks up a user by remember token and returns that user. this method
// will handle the hashing for us.
func (us *UserService) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := us.hmac.Hash(token)
	err := first(us.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// authenticate users
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}
	return foundUser, nil
}

// closes the UserService database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// ONLY FOR TEST ENVIRONMENT! drop/create user table
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

//AutoMigrate will automatically try to migrate the Users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// first will query using the provided gorm.DB and get first item returned
// and place it in destination. if nothing found by query, return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
