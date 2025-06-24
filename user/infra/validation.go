package infra

import (
	"benthos/user/dom"
	"benthos/common/res"
	"strings"
	"unicode/utf8"
)

// Validador para User
type UserValidator struct{}

func (uv *UserValidator) ValidateUser(user *dom.User) []common.Error {
    var errors []common.Error
    
    // Validar username
    if strings.TrimSpace(user.Username) == "" {
        errors = append(errors, common.Error{Message: "username required", Code: "EUSRU1"})
    } else if len(user.Username) < 6 {
        errors = append(errors, common.Error{Message: "username must have more than 6 characters", Code: "EUSRU2"})
    } else if len(user.Username) > 50 {
        errors = append(errors, common.Error{Message: "username can't have more than 50 characters", Code: "EUSRU3"})
    } else if !utf8.ValidString(user.Username) {
        errors = append(errors, common.Error{Message: "username contains invalid characters", Code: "EUSRU4"})
    }
    
    // Validar password
    if strings.TrimSpace(user.Password) == "" {
        errors = append(errors, common.Error{Message: "password required", Code: "EUSRP1"})
    } else if len(user.Password) < 6 {
        errors = append(errors, common.Error{Message: "password must have more than 6 characters", Code: "EUSRP2"})
    } else if len(user.Password) > 50 {
        errors = append(errors, common.Error{Message: "password can't have more than 50 characters", Code: "EUSRP3"})
    }
    
    return errors
}