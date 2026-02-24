package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vo1dFl0w/department-management-service/internal/domain"
	"github.com/vo1dFl0w/department-management-service/internal/repository"
)

type DepartmentManagerService interface {
	CreateNewDepartment(ctx context.Context, name string, parentID *int64) (*domain.Department, error)
	CreateNewEmployee(ctx context.Context, departmentID int64, fullName string, position string, hiredAt *time.Time) (*domain.Employee, error)
	GetDepartmentByID(ctx context.Context, id int64, depth int, includeEmployees bool) (*domain.DepartmentTree, error)
	UpdateDepartmentByID(ctx context.Context, id int64, namePtr *string, parentID *int64) (*domain.Department, error)
	DeleteDepartmentByID(ctx context.Context, id int64, mode string, reassignTo *int64) error
}

type service struct {
	storage repository.Postgres
}

func NewService(storage repository.Postgres) *service {
	return &service{storage: storage}
}

func (s *service) CreateNewDepartment(ctx context.Context, name string, parentID *int64) (*domain.Department, error) {
	if err := validateDepartmentName(name); err != nil {
		return nil, err
	}

	res, err := s.storage.CreateDepartment(ctx, name, parentID)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, domain.ErrConflict
		} else if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.ErrNotFound
		} else {
			return nil, fmt.Errorf("create department: %w", err)
		}
	}
	return res, nil
}

func (s *service) CreateNewEmployee(ctx context.Context, departmentID int64, fullName string, position string, hiredAt *time.Time) (*domain.Employee, error) {
	if err := validateFullName(fullName); err != nil {
		return nil, err
	}

	if err := validatePositionName(position); err != nil {
		return nil, err
	}

	res, err := s.storage.CreateEmployee(ctx, departmentID, fullName, position, hiredAt)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.ErrNotFound
		} else {
			return nil, fmt.Errorf("create employee: %w", err)
		}
	}
	return res, nil
}

func (s *service) GetDepartmentByID(ctx context.Context, id int64, depth int, includeEmployees bool) (*domain.DepartmentTree, error) {
	if depth == 0 {
		depth = 1
	}

	if depth > 5 {
		depth = 5
	}

	res, err := s.storage.GetDepartment(ctx, id, depth, includeEmployees)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.ErrNotFound
		} else {
			return nil, fmt.Errorf("get department: %w", err)
		}
	}
	return res, nil
}

func (s *service) UpdateDepartmentByID(ctx context.Context, id int64, namePtr *string, parentID *int64) (*domain.Department, error) {
	if parentID != nil && id == *parentID {
		return nil, domain.ErrIDMustNotReferToItself
	}

	var name string
	if namePtr != nil {
		name = strings.TrimSpace(*namePtr)
		if name != "" {
			*namePtr = name
		}
	}

	if *namePtr == "" || name == "" || namePtr == nil {
		namePtr = nil
	}

	if rs := []rune(name); len(rs) > 200 {
		return nil, domain.ErrToLongDepartmentName
	}

	res, err := s.storage.UpdateDepartment(ctx, id, namePtr, parentID)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, domain.ErrConflict
		} else if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.ErrNotFound
		} else {
			return nil, fmt.Errorf("update department: %w", err)
		}
	}
	return res, nil
}

func (s *service) DeleteDepartmentByID(ctx context.Context, id int64, mode string, reassignTo *int64) error {
	if mode == "" {
		mode = "cascade"
	}

	if mode == "reassign" && reassignTo == nil {
		return domain.ErrResignIDMustBeSet
	}

	if reassignTo != nil && *reassignTo == id {
		return domain.ErrIDMustNotReferToItself
	}

	if err := s.storage.DeleteDepartment(ctx, id, mode, reassignTo); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.ErrNotFound
		} else if errors.Is(err, repository.ErrResignIDMustBeSet) {
			return domain.ErrResignIDMustBeSet
		} else if errors.Is(err, repository.ErrIDMustNotReferToItself) {
			return domain.ErrIDMustNotReferToItself
		} else if errors.Is(err, repository.ErrConflict) {
			return domain.ErrConflict
		}
	}

	return nil
}
