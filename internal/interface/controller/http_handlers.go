package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"

	dto "github.com/fedo3nik/GamePortal_IdentityService/internal/interface/controller/dto"

	"github.com/fedo3nik/GamePortal_IdentityService/internal/application/service"
	"github.com/fedo3nik/GamePortal_IdentityService/internal/domain/entities"
	e "github.com/fedo3nik/GamePortal_IdentityService/internal/error"
	"github.com/pkg/errors"
)

type HTTPSignUpHandler struct {
	usersService service.User
}

type HTTPSignInHandler struct {
	usersService service.User
}

type HTTPAddWarnHandler struct {
	userService service.User
}

type HTTPRemWarnHandler struct {
	userService service.User
}

type HTTPIsBanedHandler struct {
	userService service.User
}

func errorType(w http.ResponseWriter, err error) {
	if errors.Is(err, e.ErrValidation) {
		_, printErr := fmt.Fprintf(w, "Request error: %v", err)
		if printErr != nil {
			log.Printf("Fprint error: %v", printErr)
		}

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if errors.Is(err, e.ErrSignIn) {
		_, printErr := fmt.Fprintf(w, "Request error: %v", err)
		if printErr != nil {
			log.Printf("Fprint error: %v", printErr)
			return
		}

		w.WriteHeader(http.StatusNotFound)

		return
	}

	if errors.Is(err, e.ErrSignUp) {
		_, printErr := fmt.Fprintf(w, "Request error: %v", err)
		if printErr != nil {
			log.Printf("Fprint error: %v", printErr)
			return
		}

		w.WriteHeader(http.StatusConflict)

		return
	}

	_, printErr := fmt.Fprintf(w, "Internal server error: %v", err)
	if printErr != nil {
		log.Printf("Fprint error: %v", printErr)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}

func NewHTTPSignInHandler(userService service.User) *HTTPSignInHandler {
	return &HTTPSignInHandler{usersService: userService}
}

func (hh HTTPSignInHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var u entities.User

	var resp dto.SignInResponse

	reqErr := json.NewDecoder(r.Body).Decode(&u)
	if reqErr != nil {
		log.Printf("Body read error: %v", reqErr)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	usr, aToken, rToken, err := hh.usersService.SignIn(r.Context(), u.Email, u.Password)
	if err != nil {
		errorType(w, err)
		return
	}

	resp.ID = usr.ID
	resp.Email = usr.Email
	resp.Nickname = usr.Nickname
	resp.AccessToken = aToken
	resp.RefreshToken = rToken

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Printf("Encode error: %v", err)
		return
	}
}

func NewHTTPSignUpHandler(usersService service.User) *HTTPSignUpHandler {
	return &HTTPSignUpHandler{usersService: usersService}
}

func (hh HTTPSignUpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var u entities.User

	var resp dto.SignUpResponse

	reqErr := json.NewDecoder(r.Body).Decode(&u)
	if reqErr != nil {
		log.Printf("Body read error: %v", reqErr)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	usr, err := hh.usersService.SignUp(r.Context(), u.Nickname, u.Email, u.Password)
	if err != nil {
		errorType(w, err)
		return
	}

	resp.ID = usr.ID
	resp.Email = usr.Email
	resp.Password = usr.Password
	resp.Nickname = usr.Nickname

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Printf("Encode error: %v", err)
		return
	}
}

func NewHTTPAddWarnHandler(userService service.User) *HTTPAddWarnHandler {
	return &HTTPAddWarnHandler{userService: userService}
}

func (hh HTTPAddWarnHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp dto.AddWarnResponse

	url := r.URL.Path
	idString := path.Base(url)

	usr, err := hh.userService.AddWarning(r.Context(), idString)
	if err != nil {
		errorType(w, err)
		return
	}

	resp.ID = usr.ID
	resp.SuccessAdd = true

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Printf("Encode error: %v", err)
		return
	}
}

func NewHTTPRemHandler(userService service.User) *HTTPRemWarnHandler {
	return &HTTPRemWarnHandler{userService: userService}
}

func (hh HTTPRemWarnHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp dto.RemWarnResponse

	url := r.URL.Path
	idString := path.Base(url)

	usr, err := hh.userService.RemoveWarning(r.Context(), idString)
	if err != nil {
		errorType(w, err)
		return
	}

	resp.ID = usr.ID
	resp.SuccessRem = true

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Printf("Encode error: %v", err)
		return
	}
}

func NewHTTPIsBannedHandler(userService service.User) *HTTPIsBanedHandler {
	return &HTTPIsBanedHandler{userService: userService}
}

func (hh HTTPIsBanedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp dto.IsBannedResponse

	url := r.URL.Path
	idString := path.Base(url)

	usr, isBanned, err := hh.userService.IsBanned(r.Context(), idString)
	if err != nil {
		errorType(w, err)
		return
	}

	resp.ID = usr.ID
	resp.IsBanned = isBanned

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Printf("Encode error: %v", err)
		return
	}
}
