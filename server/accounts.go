package server

import (
	"net/http"

	api "git.underland.io/ehazlett/fynca/api/services/accounts/v1"
	"github.com/sirupsen/logrus"
)

func (s *Server) accountsProfileHandler(w http.ResponseWriter, r *http.Request) {
	fc, err := s.getClient(r)
	if err != nil {
		logrus.WithError(err).Error("error getting fynca client")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fc.Close()

	ctx := r.Context()
	resp, err := fc.GetAccountProfile(ctx, &api.GetAccountProfileRequest{})
	if err != nil {
		logrus.WithError(err).Error("error getting account profile")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.marshaler().Marshal(w, resp); err != nil {
		logrus.WithError(err).Error("error marshaling account profile")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) accountsChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	fc, err := s.getClient(r)
	if err != nil {
		logrus.WithError(err).Error("error getting fynca client")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fc.Close()

	ctx := r.Context()

	username := r.FormValue("username")
	password := r.FormValue("password")
	currentPassword := r.FormValue("currentPassword")

	// authenticate user
	if _, err := fc.Authenticate(ctx, &api.AuthenticateRequest{
		Username: username,
		Password: []byte(currentPassword),
	}); err != nil {
		logrus.WithError(err).Error("invalid username or password")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// change
	if _, err := fc.ChangePassword(ctx, &api.ChangePasswordRequest{
		Username: username,
		Password: []byte(password),
	}); err != nil {
		logrus.WithError(err).Error("error changing password")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) accountsUpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	fc, err := s.getClient(r)
	if err != nil {
		logrus.WithError(err).Error("error getting fynca client")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fc.Close()

	ctx := r.Context()

	username := r.FormValue("username")
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	email := r.FormValue("email")

	account := &api.Account{
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
	// update
	if _, err := fc.UpdateAccount(ctx, &api.UpdateAccountRequest{
		Account: account,
	}); err != nil {
		logrus.WithError(err).Error("error updating account")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
