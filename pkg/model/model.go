package model

import (
	"time"
)

type Movie struct {
	ID        string
	MovieId   int
	Name      string
	OverView  string
	Url       string
	ImageUrl  string
	LeadActor string
	Tags      []string
	CreateTs  time.Time
	UpdateTs  time.Time
}

type DateOfBirth struct {
	Year      int
	Month     int
	MonthName string
	Date      int
}

type User struct {
	Email             string
	Username          string
	PasswordHash      string
	FullName          string
	Role              string
	Dob               DateOfBirth
	CreateTs          time.Time
	UpdateTs          time.Time
	LoginTs           time.Time
	PreviousPasswords []string
	Otp               string
}
