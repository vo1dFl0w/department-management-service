package http

import (
	"errors"
	"net/http"

	"github.com/vo1dFl0w/department-management-service/internal/domain"
	"github.com/vo1dFl0w/department-management-service/internal/transport/http/httpgen"
)

var (
	ErrBadRequest          = errors.New("bad request")
	ErrClientClosedRequest = errors.New("client closed request")
	ErrGatewayTimeout      = errors.New("gateway timeout")
	ErrInternalServerError = errors.New("internal server error")
	ErrConflict            = errors.New("conflict")
	ErrNotFound            = errors.New("not found")
)

type HTTPError struct {
	Message string
	Status  int
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) ToDepartmentsPostRes() httpgen.DepartmentsPostRes {
	switch e.Status {
	case http.StatusBadRequest:
		return &httpgen.DepartmentsPostBadRequest{Message: e.Message, Status: e.Status}
	case http.StatusNotFound:
		return &httpgen.DepartmentsPostNotFound{Message: e.Message, Status: e.Status}
	case http.StatusConflict:
		return &httpgen.DepartmentsPostConflict{Message: e.Message, Status: e.Status}
	case http.StatusGatewayTimeout:
		return &httpgen.DepartmentsPostGatewayTimeout{Message: e.Message, Status: e.Status}
	default:
		return &httpgen.DepartmentsPostInternalServerError{Message: e.Message, Status: e.Status}
	}
}

func (e *HTTPError) ToDepartmentsIDPatchRes() httpgen.DepartmentsIDPatchRes {
	switch e.Status {
	case http.StatusBadRequest:
		return &httpgen.DepartmentsIDPatchBadRequest{Message: e.Message, Status: e.Status}
	case http.StatusNotFound:
		return &httpgen.DepartmentsIDPatchNotFound{Message: e.Message, Status: e.Status}
	case http.StatusConflict:
		return &httpgen.DepartmentsIDPatchConflict{Message: e.Message, Status: e.Status}
	case http.StatusGatewayTimeout:
		return &httpgen.DepartmentsIDPatchGatewayTimeout{Message: e.Message, Status: e.Status}
	default:
		return &httpgen.DepartmentsIDPatchInternalServerError{Message: e.Message, Status: e.Status}
	}
}

func (e *HTTPError) ToDepartmentsIDGetRes() httpgen.DepartmentsIDGetRes {
	switch e.Status {
	case http.StatusBadRequest:
		return &httpgen.DepartmentsIDGetBadRequest{Message: e.Message, Status: e.Status}
	case http.StatusNotFound:
		return &httpgen.DepartmentsIDGetNotFound{Message: e.Message, Status: e.Status}
	case http.StatusGatewayTimeout:
		return &httpgen.DepartmentsIDGetGatewayTimeout{Message: e.Message, Status: e.Status}
	default:
		return &httpgen.DepartmentsIDGetInternalServerError{Message: e.Message, Status: e.Status}
	}
}

func (e *HTTPError) ToDepartmentsIDEmployeesPostRes() httpgen.DepartmentsIDEmployeesPostRes {
	switch e.Status {
	case http.StatusBadRequest:
		return &httpgen.DepartmentsIDEmployeesPostBadRequest{Message: e.Message, Status: e.Status}
	case http.StatusNotFound:
		return &httpgen.DepartmentsIDEmployeesPostNotFound{Message: e.Message, Status: e.Status}
	case http.StatusGatewayTimeout:
		return &httpgen.DepartmentsIDEmployeesPostGatewayTimeout{Message: e.Message, Status: e.Status}
	default:
		return &httpgen.DepartmentsIDEmployeesPostInternalServerError{Message: e.Message, Status: e.Status}
	}
}

func (e *HTTPError) ToDepartmentsIDDeleteRes() httpgen.DepartmentsIDDeleteRes {
	switch e.Status {
	case http.StatusBadRequest:
		return &httpgen.DepartmentsIDDeleteBadRequest{Message: e.Message, Status: e.Status}
	case http.StatusNotFound:
		return &httpgen.DepartmentsIDDeleteNotFound{Message: e.Message, Status: e.Status}
	case http.StatusGatewayTimeout:
		return &httpgen.DepartmentsIDDeleteGatewayTimeout{Message: e.Message, Status: e.Status}
	default:
		return &httpgen.DepartmentsIDDeleteInternalServerError{Message: e.Message, Status: e.Status}
	}
}

func MapError(err error) *HTTPError {
	switch {
	case errors.Is(err, domain.ErrEmptyDepartmentName) || errors.Is(err, domain.ErrEmptyFullName) || errors.Is(err, domain.ErrEmptyPositionName) ||
		errors.Is(err, domain.ErrIDMustNotReferToItself) || errors.Is(err, domain.ErrResignIDMustBeSet) || errors.Is(err, domain.ErrToLongDepartmentName) ||
		errors.Is(err, domain.ErrToLongFullName) || errors.Is(err, domain.ErrToLongPositionName):
		return &HTTPError{Message: ErrBadRequest.Error(), Status: http.StatusBadRequest}
	case errors.Is(err, domain.ErrNotFound):
		return &HTTPError{Message: ErrNotFound.Error(), Status: http.StatusNotFound}
	case errors.Is(err, domain.ErrConflict):
		return &HTTPError{Message: ErrConflict.Error(), Status: http.StatusConflict}
	case errors.Is(err, domain.ErrGatewayTimeout):
		return &HTTPError{Message: ErrGatewayTimeout.Error(), Status: http.StatusGatewayTimeout}
	default:
		return &HTTPError{Message: ErrInternalServerError.Error(), Status: http.StatusInternalServerError}
	}
}
