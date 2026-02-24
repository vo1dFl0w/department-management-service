package repository

import (
	"context"
	"time"

	"github.com/vo1dFl0w/department-management-service/internal/domain"
)

type Postgres interface {
	CreateDepartment(ctx context.Context, name string, parentID *int64) (*domain.Department, error)
	CreateEmployee(ctx context.Context, departmentID int64, fullName string, position string, hiredAt *time.Time) (*domain.Employee, error)
	GetDepartment(ctx context.Context, id int64, depth int, includeEmployees bool) (*domain.DepartmentTree, error)
	UpdateDepartment(ctx context.Context, id int64, namePtr *string, parentID *int64) (*domain.Department, error)
	DeleteDepartment(ctx context.Context, id int64, mode string, reassignTo *int64) error
}