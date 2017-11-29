package models

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"../hash"
	"../rand"
)

var (
	// ErrNotFound returned when resource cannot be found in the database
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalid returned when an invalid ID is provided to method like Delete
	ErrInvalidID = errors.New("models: ID must be > 0")

	// ErrInvalidPassword returned when incorrect password used for authentication
	ErrInvalidPassword = errors.New("models: incorrect password provided")

	// ErrEmailRequired returned when is email is blank on form submission
	ErrEmailRequired = errors.New("modesl: Email address is required")

	// ErrEmailInvalid returned when email syntax is bad
	ErrEmailInvalid = errors.New("models: Email address is not valid")

	// emailRegex is used to validate email addresses
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

const userPwPepper = "wtul91.5"
const hmacSecretKey = "progressivestereo"

// User represents the user model stored in the database.
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `sql:"-"` //telling gorm to ignore this - don't store in db
	PasswordHash string `gorm:"not null"`
	Remember     string `sql:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// UserDB is used to interact with the Users database
type UserDB interface {
	// Methods for querying single users
	ById(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Close closes a DB connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

// UserService is a set of methods used to work with the user model
type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

// look up user by id.
type userService struct {
	UserDB
}

// database connection
type userGorm struct {
	db *gorm.DB
}

// userValidator is layer that validates data before going to DB
type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

// test that userGorm implements UserDB interface
var _ UserDB = &userGorm{}

// test that userValidator implements UserDB interface
var _ UserDB = &userValidator{}

//test that userService implements UserDB interface
var _ UserDB = &userService{}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

// create userGorm -- establish database connection
func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &userGorm{
		db: db,
	}, nil
}

func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := &userValidator{
		hmac:   hmac,
		UserDB: ug,
	}

	return &userService{
		UserDB: uv,
	}, nil
}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: emailRegex,
	}
}

// ByRemember will hash the remember token and then call
// ByRemember on the subsequent DB layer
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

// create user
func (uv *userValidator) Create(user *User) error {
	err := runUserValFuncs(user,
		uv.bcryptPassword,
		uv.setRememberIfUnset,
		uv.hmacRemember,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

// Create will create the provided user
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Update will hash a remember token if it is provided
func (uv *userValidator) Update(user *User) error {
	err := runUserValFuncs(user,
		uv.bcryptPassword,
		uv.hmacRemember,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// update user (expects complete User object)
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFuncs(&user, uv.idNotZero)
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

// Delete will remove user with provided ID from the DB
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// ById searches for user by id
func (ug *userGorm) ById(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail normalizes the email before calling ByEmail on
// the UserDB layer
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	if err := runUserValFuncs(&user, uv.normalizeEmail); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

// look up user by email
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// looks up a user by remember token and returns that user. this method
// expects the remember token to already be hashed.
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// authenticate users
func (us *userService) Authenticate(email, password string) (*User, error) {
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

// bcryptPassword will hash a user's password with a
// predefined pepper (userPwPepper) and bcrypt
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(pwBytes), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember == "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func (uv *userValidator) idNotZero(user *User) error {
	if user.ID == 0 {
		return ErrInvalidID
	}
	return nil
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Password) {
		return ErrEmailInvalid
	}
	return nil
}

// closes the UserService database connection
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// ONLY FOR TEST ENVIRONMENT! drop/create user table
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

//AutoMigrate will automatically try to migrate the Users table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
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
