package model

import "time"

type PostUsersInput struct {
	Age  uint16
	Name string
}

type UserOutput struct {
	Id   uint64
	Age  uint16
	Name string
	Time time.Time
}
