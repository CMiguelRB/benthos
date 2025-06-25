package infra

import (
	commonDom "benthos/common/dom"
	userDom "benthos/user/dom"
	"strings"
	"unicode/utf8"
)

type UserValidator struct{}

func (uv *UserValidator) ValidateUser(user *userDom.User) []commonDom.Error {
	var errors []commonDom.Error

	if strings.TrimSpace(user.Username) == "" {
		errors = append(errors, commonDom.Error{Message: "username required", Code: "EUSRU1"})
	} else if len(user.Username) < 6 {
		errors = append(errors, commonDom.Error{Message: "username must have more than 6 characters", Code: "EUSRU2"})
	} else if len(user.Username) > 50 {
		errors = append(errors, commonDom.Error{Message: "username can't have more than 50 characters", Code: "EUSRU3"})
	} else if !utf8.ValidString(user.Username) {
		errors = append(errors, commonDom.Error{Message: "username contains invalid characters", Code: "EUSRU4"})
	}

	if strings.TrimSpace(user.Password) == "" {
		errors = append(errors, commonDom.Error{Message: "password required", Code: "EUSRP1"})
	} else if len(user.Password) < 6 {
		errors = append(errors, commonDom.Error{Message: "password must have more than 6 characters", Code: "EUSRP2"})
	} else if len(user.Password) > 50 {
		errors = append(errors, commonDom.Error{Message: "password can't have more than 50 characters", Code: "EUSRP3"})
	}

	return errors
}
