package common

import "fmt"

const (
	DBTypeUser  = 2
	CurrentUser = "user"
)

const DateTimeFmt = "2006-01-02 15:04:05.999999"

type Requester interface {
	GetUserId() int
	GetEmail() string
	GetRole() string
}

func AppRecover() {
	if err := recover(); err != nil {
		fmt.Println("Recover: ", err)
	}
}
