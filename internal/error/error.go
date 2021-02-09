package error

import "github.com/pkg/errors"

var ErrValidation = errors.New("validation error")
var ErrSignUp = errors.New("sign up error")
var ErrSignIn = errors.New("sign in error")
var ErrDB = errors.New("database error")
