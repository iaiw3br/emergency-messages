package controllers

import (
	"github.com/pkg/errors"
	"net/http"
	"projects/emergency-messages/internal/errorx"
)

func assertError(err error, w http.ResponseWriter) bool {
	if err == nil {
		return false
	}
	switch errors.Cause(err) {
	case errorx.ErrNotFound:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	case errorx.ErrValidation:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	return true
}
