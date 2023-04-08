package svc

import (
	"crypto/rand"
	"log"
)

const (
	otpChars  = "1234567890"
	otpLength = 6
)

func (s *svc) GenerateOtp() (string, error) {
	buffer := make([]byte, otpLength)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	otpCharslen := len(otpChars)
	for i := 0; i < otpLength; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharslen]
	}
	return string(buffer), nil
}

func (s *svc) VerifyOtp(email string, otp string) (bool, error) {
	if email == "" || otp == "" {
		log.Println("email or otp is empty empty field ..")
		return false, ErrBadRequest
	}
	storedOtp, err := s.cache.GetFromCache(email)
	if err != nil {
		log.Println("error in getting otp from cache...", err)
		return false, err
	}
	if storedOtp == otp {
		log.Println("otp matched...")
		return true, nil
	}
	log.Println("otp mismatched...")
	return false, nil
}


