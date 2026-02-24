package http

import (
	"context"
	"net/http"

	"github.com/vo1dFl0w/department-management-service/internal/config"
	"github.com/vo1dFl0w/department-management-service/internal/transport/http/httpgen"
	"github.com/vo1dFl0w/department-management-service/internal/usecase"
	"github.com/vo1dFl0w/department-management-service/pkg/logger"
)

type Handler struct {
	cfg                  *config.Config
	logger               logger.Logger
	departmentManagerSrv usecase.DepartmentManagerService
}

func NewHandler(cfg *config.Config, logger logger.Logger, departmentManagerSrv usecase.DepartmentManagerService) *Handler {
	return &Handler{
		cfg:                  cfg,
		logger:               logger,
		departmentManagerSrv: departmentManagerSrv,
	}
}

func (h *Handler) DepartmentsPost(ctx context.Context, req *httpgen.CreateDepartmentRequest) (httpgen.DepartmentsPostRes, error) {
	parentID := optNilIntToInt64Ptr(req.ParentID)
	dep, err := h.departmentManagerSrv.CreateNewDepartment(ctx, req.Name, parentID)
	if err != nil {
		httpErr := MapError(err)
		h.LogHTTPError(ctx, err, httpErr)
		return httpErr.ToDepartmentsPostRes(), nil
	}
	res := domainDepartmentToHTTP(*dep)

	return &res, nil
}

func (h *Handler) DepartmentsIDPatch(ctx context.Context, req *httpgen.UpdateDepartmentRequest, params httpgen.DepartmentsIDPatchParams) (httpgen.DepartmentsIDPatchRes, error) {
	parentID := int64(req.ParentID.Value)

	dep, err := h.departmentManagerSrv.UpdateDepartmentByID(ctx, int64(params.ID), &req.Name.Value, &parentID)
	if err != nil {
		httpErr := MapError(err)
		h.LogHTTPError(ctx, err, httpErr)
		return httpErr.ToDepartmentsIDPatchRes(), nil
	}

	res := domainDepartmentToHTTP(*dep)

	return &res, nil
}

func (h *Handler) DepartmentsIDGet(ctx context.Context, params httpgen.DepartmentsIDGetParams) (httpgen.DepartmentsIDGetRes, error) {
	depTree, err := h.departmentManagerSrv.GetDepartmentByID(ctx, int64(params.ID), params.Depth.Value, params.IncludeEmployees.Value)
	if err != nil {
		httpErr := MapError(err)
		h.LogHTTPError(ctx, err, httpErr)
		return httpErr.ToDepartmentsIDGetRes(), nil
	}

	res := domainTreeToHTTP(depTree)

	return &res, nil
}

func (h *Handler) DepartmentsIDEmployeesPost(ctx context.Context, req *httpgen.CreateEmployeeRequest, params httpgen.DepartmentsIDEmployeesPostParams) (httpgen.DepartmentsIDEmployeesPostRes, error) {
	emp, err := h.departmentManagerSrv.CreateNewEmployee(ctx, int64(params.ID), req.FullName, req.Position, &req.HiredAt.Value)
	if err != nil {
		httpErr := MapError(err)
		h.LogHTTPError(ctx, err, httpErr)
		return httpErr.ToDepartmentsIDEmployeesPostRes(), nil
	}

	res := domainEmployeeToHTTP(*emp)

	return &res, nil
}

func (h *Handler) DepartmentsIDDelete(ctx context.Context, params httpgen.DepartmentsIDDeleteParams) (httpgen.DepartmentsIDDeleteRes, error) {
	ressignTo := int64(params.ReassignToDepartmentID.Value)

	if err := h.departmentManagerSrv.DeleteDepartmentByID(ctx, int64(params.ID), string(params.Mode), &ressignTo); err != nil {
		httpErr := MapError(err)
		h.LogHTTPError(ctx, err, httpErr)
		return httpErr.ToDepartmentsIDDeleteRes(), nil
	}

	return &httpgen.DepartmentsIDDeleteNoContent{}, nil
}

func (h *Handler) LogHTTPError(ctx context.Context, err error, httpErr *HTTPError) {
	attrs := []any{
		"error", err,
		"status", httpErr.Status,
		"message", httpErr.Message,
	}

	switch {
	case httpErr.Status >= 500:
		switch httpErr.Status {
		case http.StatusGatewayTimeout:
			h.logger.Error("http_request_failed", append(attrs, "reason", "dependency_timeout")...)
		default:
			h.logger.Error("http_request_failed", append(attrs, "reason", "internal_server_error")...)
		}
	case httpErr.Status >= 400:
		h.logger.Warn("http_request_failed", append(attrs, "reason", "client_error")...)
	}
}
