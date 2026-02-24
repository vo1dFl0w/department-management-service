package postgres

import "time"

type Department struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	ParentID  *int64    `gorm:"column:parent_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

type Employee struct {
	ID           int64      `gorm:"column:id;primaryKey"`
	DepartmentID int64      `gorm:"column:department_id"`
	FullName     string     `gorm:"column:full_name"`
	Position     string     `gorm:"column:position"`
	HiredAt      *time.Time `gorm:"column:hired_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
}

type DepartmentTree struct {
	Department Department        `json:"department"`
	Employees  []Employee        `json:"employees"`
	Children   []*DepartmentTree `json:"children"`
}
