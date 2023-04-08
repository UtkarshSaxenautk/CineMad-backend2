package svc

import (
	"authentication-ms/pkg/model"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type svc struct {
	dao   Dao
	sdk   Sdk
	cache Cache
}

var (
	//these are public errors, please do not include any technical wordings here
	ErrNoData               = errors.New("no data found")
	ErrBadRequest           = errors.New("mandatory input data missing")
	ErrDeleteFailed         = errors.New("domain delete failed")
	ErrUnexpected           = errors.New("unexpected error")
	ErrEmailAlreadyInUse    = errors.New("email is already in use")
	ErrUserNameAlreadyInUse = errors.New("username is already in use")
	ErrUserNotAuthorized    = errors.New("user not authenticated")
)

func New(dao Dao, cache Cache, sdk Sdk) SVC {
	s := &svc{dao, sdk, cache}
	return s
}

func (s *svc) Signup(ctx context.Context, user model.User) error {
	emptyDob := model.DateOfBirth{
		Date:      -1,
		Year:      -1,
		Month:     -1,
		MonthName: "",
	}
	log.Println(user)
	if user.Email == "" || user.Username == "" || user.PasswordHash == "" || user.Dob == emptyDob || user.FullName == "" {
		log.Println("missing necessary field...")
		return ErrBadRequest
	}
	emailExist, userNameExist, err := s.dao.CheckEmailAndUserName(ctx, user)
	if err != nil {
		log.Println("error in checking email and username existence..")
		return err
	}
	if emailExist {
		log.Println("email already in use...")
		return ErrEmailAlreadyInUse
	}
	if userNameExist {
		log.Println("username already in use...")
		return ErrUserNameAlreadyInUse
	}
	user.PasswordHash, err = s.hashPassword(user.PasswordHash)
	if err != nil {
		log.Println("error in creating password hash: ", err)
		return ErrUnexpected
	}
	err = s.dao.CreateUser(ctx, user)
	if err != nil {
		log.Println("error in creating user...")
		return ErrUnexpected
	}
	log.Println("user created successfully")
	return nil
}

func (s *svc) ForgotPassword(email string) error {
	return nil
}

func (s *svc) hashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with Bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

func (s *svc) passwordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

func (s *svc) SignIn(ctx context.Context, email string, password string) (bool, error) {
	if email == "" || password == "" {
		return false, ErrBadRequest
	}
	hashedPassword, err := s.dao.GetUser(ctx, email)

	if err != nil {
		log.Println("error in getting user from email")
		return false, ErrUnexpected
	}

	matched := s.passwordsMatch(hashedPassword, password)
	if matched {
		log.Println("password matched...")
		return true, nil
	}
	log.Println("password mismatched...")
	return false, ErrUserNotAuthorized
}

func (s *svc) ChangePassword(ctx context.Context, user model.User, newPassword string) error {
	if user.Email == "" || user.PasswordHash == "" || newPassword == "" {
		log.Println("necessary field in missing...")
		return ErrBadRequest
	}
	log.Println("start.....")
	newHashed, err := s.hashPassword(newPassword)
	if err != nil {
		log.Println("error in hashing new password..", err)
		return ErrUnexpected
	}
	log.Println("new password hashed")
	err = s.dao.UpdatePassword(ctx, user, newHashed)
	if err != nil {
		log.Println("error in updating password ..", err)
		return ErrUnexpected
	}
	log.Println("password successfully changed...")
	return nil
}
