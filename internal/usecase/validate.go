package usecase

import "github.com/vo1dFl0w/department-management-service/internal/domain"

func validateDepartmentName(s string) error {
	if s == "" {
		return domain.ErrEmptyDepartmentName
	}

	rs := []rune(s)
	if len(rs) > 200 {
		return domain.ErrToLongDepartmentName
	}
	return nil
}

func validateFullName(s string) error {
	if s == "" {
		return domain.ErrEmptyFullName
	}

	rs := []rune(s)
	if len(rs) > 200 {
		return domain.ErrToLongFullName
	}
	return nil
}

func validatePositionName(s string) error {
	if s == "" {
		return domain.ErrEmptyPositionName
	}

	rs := []rune(s)
	if len(rs) > 200 {
		return domain.ErrToLongPositionName
	}
	return nil
}
