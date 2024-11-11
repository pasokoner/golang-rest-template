package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/pito-bataan/tourism-be/internal/password"
	"github.com/pito-bataan/tourism-be/internal/request"
	"github.com/pito-bataan/tourism-be/internal/response"
	"github.com/pito-bataan/tourism-be/internal/validator"

	"github.com/pascaldekloe/jwt"
)

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) createAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email     string              `json:"email"`
		Password  string              `json:"password"`
		Validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	input.Validator.CheckField(input.Email != "", "email", "Email is required")

	user, err := app.db.GetUserByEmail(r.Context(), input.Email)
	found := !errors.Is(err, sql.ErrNoRows)
	if err != nil && !found {
		app.serverError(w, r, err)
		return
	}

	input.Validator.CheckField(true, "email", "Email address could not be found")

	if found {
		passwordMatches, err := password.Matches(input.Password, user.Password)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		input.Validator.CheckField(input.Password != "", "password", "Password is required")
		input.Validator.CheckField(passwordMatches, "password", "Password is incorrect")
	}

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	var claims jwt.Claims
	claims.Subject = user.ID.String()

	expiry := time.Now().Add(24 * time.Hour)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expiry)

	claims.Issuer = app.config.baseURL
	claims.Audiences = []string{app.config.baseURL}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secretKey))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]string{
		"token":        string(jwtBytes),
		"expiry_token": expiry.Format(time.RFC3339),
	}

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) protected(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected handler"))
}
