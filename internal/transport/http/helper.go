package http

import (
	"time"

	"github.com/vo1dFl0w/department-management-service/internal/domain"
	"github.com/vo1dFl0w/department-management-service/internal/transport/http/httpgen"
)

func optNilIntToInt64Ptr(opt httpgen.OptNilInt) *int64 {
	if !opt.Set || opt.Null {
		return nil
	}
	v := int64(opt.Value)
	return &v
}

func toOptNilInt(p *int64) httpgen.OptNilInt {
	if p == nil {
		return httpgen.OptNilInt{}
	}
	return httpgen.NewOptNilInt(int(*p))
}

func toOptNilDate(t *time.Time) httpgen.OptNilDate {
	if t == nil {
		return httpgen.OptNilDate{}
	}
	return httpgen.NewOptNilDate(*t)
}

func domainEmployeeToHTTP(e domain.Employee) httpgen.Employee {
	return httpgen.Employee{
		ID:           int(e.ID),
		DepartmentID: int(e.DepartmentID),
		FullName:     e.FullName,
		Position:     e.Position,
		HiredAt:      toOptNilDate(e.HiredAt),
		CreatedAt:    e.CreatedAt,
	}
}

func domainDepartmentToHTTP(d domain.Department) httpgen.Department {
	return httpgen.Department{
		ID:        int(d.ID),
		Name:      d.Name,
		ParentID:  toOptNilInt(d.ParentID),
		CreatedAt: d.CreatedAt,
	}
}

func domainTreeToHTTP(t *domain.DepartmentTree) httpgen.DepartmentWithChildrenAndEmployees {
	emps := make([]httpgen.Employee, 0, len(t.Employees))
	for _, e := range t.Employees {
		emps = append(emps, domainEmployeeToHTTP(e))
	}

	children := make([]httpgen.DepartmentWithChildrenAndEmployees, 0, len(t.Children))
	for _, ch := range t.Children {
		children = append(children, domainTreeToHTTP(ch))
	}

	return httpgen.DepartmentWithChildrenAndEmployees{
		Department: domainDepartmentToHTTP(t.Department),
		Employees:  emps,
		Children:   children,
	}
}
