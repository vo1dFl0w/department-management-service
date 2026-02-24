package domain

import "time"

type Department struct {
	ID        int64
	Name      string
	ParentID  *int64
	CreatedAt time.Time
}

type Employee struct {
	ID           int64
	DepartmentID int64
	FullName     string
	Position     string
	HiredAt      *time.Time
	CreatedAt    time.Time
}

type DepartmentTree struct {
	Department Department
	Employees  []Employee
	Children   []*DepartmentTree
}
