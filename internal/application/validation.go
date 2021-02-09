package application

import (
	"log"
	"regexp"

	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	passwordValidator "github.com/lane-c-wagner/go-password-validator"
)

func IsEmailValid(email string) bool {
	eLen := len(email)
	if eLen < 3 || eLen > 255 {
		return false
	}

	v := ev.NewSyntaxValidator().Validate(ev.NewInput(evmail.FromString(email)))
	if !v.IsValid() {
		log.Printf("Email is invalid")
	} else {
		return true
	}

	return false
}

func IsPasswordValid(pass string) bool {
	pLen := len(pass)
	if pLen < 5 || pLen > 45 {
		return false
	}

	const minEntropyBits = 20

	err := passwordValidator.Validate(pass, minEntropyBits)
	if err != nil {
		log.Printf("Password is invalid")
	} else {
		return true
	}

	return false
}

func IsNicknameValid(nick string) bool {
	nLen := len(nick)
	if nLen < 3 || nLen > 45 {
		return false
	}

	nicknameRegex := regexp.MustCompile("[a-zA-Z0-9]{3,}")
	matched := nicknameRegex.MatchString(nick)

	return matched
}
