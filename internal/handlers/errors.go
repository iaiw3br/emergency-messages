package handlers

import (
	"github.com/pkg/errors"
	"net/http"
	"projects/emergency-messages/internal/errorx"
)

func assertError(err error) int {
	if err == nil {
		return 0
	}
	switch errors.Cause(err) {
	case errorx.ErrNotFound:
		return http.StatusNotFound
	case errorx.ErrValidation:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
