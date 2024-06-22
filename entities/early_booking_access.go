package entities

import "time"

type EBA struct {
	UserID        string
	Slot          int
	ExpiredDate   time.Time
	AvailableFrom time.Time
}

type BYPASS struct {
	Uid    string
	Bypass bool
}
