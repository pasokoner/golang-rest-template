package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/pito-bataan/tourism-be/internal/database"
	"github.com/pito-bataan/tourism-be/internal/password"
	"github.com/pito-bataan/tourism-be/internal/request"
	"github.com/pito-bataan/tourism-be/internal/response"
	"github.com/pito-bataan/tourism-be/internal/validator"
)

func (app *application) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	user := contextGetAuthenticatedUser(r)

	response.JSON(w, http.StatusOK, user)
}

func (app *application) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string              `json:"name"`
		Email     string              `json:"email"`
		Password  string              `json:"password"`
		Validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	_, err = app.db.GetUserByEmail(r.Context(), input.Email)
	found := !errors.Is(err, sql.ErrNoRows)
	if err != nil && found {
		app.serverError(w, r, err)
		return
	}

	input.Validator.CheckField(input.Email != "", "email", "Email is required")
	input.Validator.CheckField(validator.Matches(input.Email, validator.RgxEmail), "email", "Must be a valid email address")
	input.Validator.CheckField(!found, "email", "Email is already in use")

	input.Validator.CheckField(input.Password != "", "password", "Password is required")
	input.Validator.CheckField(len(input.Password) >= 8, "password", "Password is too short")
	input.Validator.CheckField(len(input.Password) <= 72, "password", "Password is too long")
	input.Validator.CheckField(validator.NotIn(input.Password, password.CommonPasswords...), "password", "Password is too common")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	_, err = app.db.CreateUser(r.Context(), database.CreateUserParams{

		Email:    input.Email,
		Password: hashedPassword,
		Name:     input.Name,
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) handlerUpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CurrentPassword string              `json:"current_password"`
		Password        string              `json:"password"`
		Validator       validator.Validator `json:"-"`
	}

	user := contextGetAuthenticatedUser(r)

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	passwordMatches, err := password.Matches(input.CurrentPassword, user.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	input.Validator.CheckField(passwordMatches, "password", "Password is incorrect")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	_, err = app.db.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		ID:       user.ID,
		Password: hashedPassword,
	})

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

func (app *application) handlerUpdateUserEmail(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email     string              `json:"email"`
		Validator validator.Validator `json:"-"`
	}

	user := contextGetAuthenticatedUser(r)

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	_, err = app.db.GetUserByEmail(r.Context(), input.Email)
	found := !errors.Is(err, sql.ErrNoRows)
	if err != nil && found {
		app.serverError(w, r, err)
		return
	}

	input.Validator.CheckField(input.Email != "", "email", "Email is required")
	input.Validator.CheckField(validator.Matches(input.Email, validator.RgxEmail), "email", "Must be a valid email address")
	input.Validator.CheckField(!found, "email", "Email is already in use")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	_, err = app.db.UpdateUserEmail(r.Context(), database.UpdateUserEmailParams{
		ID:    user.ID,
		Email: input.Email,
	})

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}
