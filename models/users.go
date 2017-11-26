package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound returned when resource cannot be found in the database
	ErrNotFound  = errors.New("models: resource not found.")
	ErrInvalidID = errors.New("models: ID must be > 0")
)

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"` //telling gorm to ignore this - don't store in db
	PasswordHash string `gorm:"not null"`
}

// look up user by id.
type UserService struct {
	db *gorm.DB
}

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
}

// create user
func (us *UserService) Create(user *User) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = "" //not entirely necessary, but could help preventing raw password from appearing in logs
	return us.db.Create(user).Error
}

// update user (expects complete User object)
func (us *UserService) Update(user *User) error {
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
