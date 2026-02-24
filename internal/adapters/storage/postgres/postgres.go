package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/vo1dFl0w/department-management-service/internal/domain"
	"github.com/vo1dFl0w/department-management-service/internal/repository"
	"gorm.io/gorm"
)

type postgres struct {
	db *gorm.DB
}

func NewPostgres(db *gorm.DB) *postgres {
	return &postgres{db: db}
}

func (p *postgres) CreateDepartment(ctx context.Context, name string, parentID *int64) (*domain.Department, error) {
	tx := p.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var cnt int64
	if err := tx.Model(&Department{}).
		Where("COALESCE(parent_id, 0) = COALESCE(?, 0) AND lower(trim(name)) = lower(trim(?))", parentID, name).
		Count(&cnt).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if cnt > 0 {
		tx.Rollback()
		return nil, repository.ErrConflict
	}

	if parentID != nil {
		var parent Department
		if err := tx.First(&parent, *parentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return nil, repository.ErrNotFound
			}
			tx.Rollback()
			return nil, err
		}
	}

	now := time.Now().UTC()
	dep := &Department{
		Name:      name,
		ParentID:  parentID,
		CreatedAt: now,
	}
	if err := tx.Create(dep).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return &domain.Department{
		ID:        dep.ID,
		Name:      name,
		ParentID:  parentID,
		CreatedAt: now,
	}, nil
}

func (p *postgres) CreateEmployee(ctx context.Context, departmentID int64, fullName string, position string, hiredAt *time.Time) (*domain.Employee, error) {
	tx := p.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var dep Department
	if err := tx.First(&dep, departmentID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	now := time.Now().UTC()
	emp := &Employee{
		DepartmentID: departmentID,
		FullName:     fullName,
		Position:     position,
		HiredAt:      hiredAt,
		CreatedAt:    now,
	}
	if err := tx.Create(emp).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return &domain.Employee{
		ID:           emp.ID,
		DepartmentID: emp.DepartmentID,
		FullName:     emp.FullName,
		Position:     emp.Position,
		HiredAt:      emp.HiredAt,
		CreatedAt:    emp.CreatedAt,
	}, nil
}

func (p *postgres) GetDepartment(ctx context.Context, id int64, depth int, includeEmployees bool) (*domain.DepartmentTree, error) {
	if depth <= 0 {
		depth = 1
	}

	if depth > 5 {
		depth = 5
	}

	type DeptRow struct {
		ID        int64
		Name      string
		ParentID  sql.NullInt64
		CreatedAt time.Time
		Lvl       int
	}

	cte := `
	WITH RECURSIVE subtree AS (
  		SELECT id, name, parent_id, created_at, 1 AS lvl
  		FROM departments
  		WHERE id = $1
  		UNION ALL
  		SELECT d.id, d.name, d.parent_id, d.created_at, s.lvl + 1
  		FROM departments d
  		JOIN subtree s ON d.parent_id = s.id
  		WHERE s.lvl < $2
	)
	SELECT id, name, parent_id, created_at, lvl
	FROM subtree
	ORDER BY lvl, parent_id NULLS FIRST, id;
	`

	var rows []DeptRow
	if err := p.db.WithContext(ctx).
		Raw(cte, id, depth).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, repository.ErrNotFound
	}

	// build nodes map (domain nodes immediately)
	nodes := make(map[int64]*domain.DepartmentTree, len(rows))
	for _, r := range rows {
		var parentPtr *int64
		if r.ParentID.Valid {
			v := r.ParentID.Int64
			parentPtr = &v
		}

		nodes[r.ID] = &domain.DepartmentTree{
			Department: domain.Department{
				ID:        r.ID,
				Name:      r.Name,
				ParentID:  parentPtr,
				CreatedAt: r.CreatedAt,
			},
			Employees: []domain.Employee{},
			Children:  []*domain.DepartmentTree{},
		}
	}

	if includeEmployees {
		ids := make([]int64, 0, len(nodes))
		for did := range nodes {
			ids = append(ids, did)
		}

		if len(ids) > 0 {
			var employees []domain.Employee
			if err := p.db.WithContext(ctx).
				Where("department_id IN ?", ids).
				Order("created_at, full_name").
				Find(&employees).Error; err != nil {
				return nil, err
			}

			for _, e := range employees {
				if node, ok := nodes[e.DepartmentID]; ok {
					node.Employees = append(node.Employees, e)
				}
			}
		}
	}

	for _, r := range rows {
		node := nodes[r.ID]

		if r.ParentID.Valid {
			parentID := r.ParentID.Int64

			if parentNode, exists := nodes[parentID]; exists {
				parentNode.Children = append(parentNode.Children, node)
			}
		}
	}

	rootNode := nodes[id]
	if rootNode == nil {
		return nil, repository.ErrNotFound
	}

	return rootNode, nil
}

func (p *postgres) UpdateDepartment(ctx context.Context, id int64, namePtr *string, parentID *int64) (*domain.Department, error) {
	tx := p.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var dep Department
	if err := tx.Clauses().First(&dep, id).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	if parentID != nil {
		var pdep Department
		if err := tx.First(&pdep, *parentID).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, repository.ErrNotFound
			}
			return nil, err
		}

		exists, err := p.checkParentInSubtreeTx(ctx, tx, id, *parentID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if exists {
			tx.Rollback()
			return nil, repository.ErrNotFound
		}
	}

	targetParent := dep.ParentID
	if parentID != nil {
		targetParent = parentID
	}
	if namePtr != nil {
		var cnt int64
		if err := tx.Model(&Department{}).
			Where("COALESCE(parent_id,0) = COALESCE(?,0) AND lower(trim(name)) = lower(trim(?)) AND id <> ?", targetParent, *namePtr, id).
			Count(&cnt).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if cnt > 0 {
			tx.Rollback()
			return nil, repository.ErrConflict
		}
	}

	updates := map[string]interface{}{}
	if namePtr != nil {
		updates["name"] = *namePtr
	}
	if parentID != nil {
		updates["parent_id"] = parentID
	}
	if len(updates) > 0 {
		if err := tx.Model(&Department{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	var updated Department
	if err := tx.First(&updated, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return &domain.Department{
		ID: updated.ID,
		Name: updated.Name,
		ParentID: updated.ParentID,
		CreatedAt: updated.CreatedAt,
	}, nil
}

func (p *postgres) checkParentInSubtreeTx(ctx context.Context, tx *gorm.DB, deptID, candidateParentID int64) (bool, error) {
	var exists int
	cte := `
	WITH RECURSIVE subtree AS (
	  SELECT id FROM departments WHERE id = ?
	  UNION ALL
	  SELECT d.id FROM departments d JOIN subtree s ON d.parent_id = s.id
	)
	SELECT 1 FROM subtree WHERE id = ? LIMIT 1;
	`
	row := tx.Raw(cte, deptID, candidateParentID).Row()
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return exists == 1, nil
}

func (p *postgres) DeleteDepartment(ctx context.Context, id int64, mode string, reassignTo *int64) error {
	tx := p.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var dep Department
	if err := tx.First(&dep, id).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.ErrNotFound
		}
		return err
	}

	if mode == "reassign" {
		if reassignTo == nil {
			tx.Rollback()
			return repository.ErrResignIDMustBeSet
		}
		if *reassignTo == id {
			tx.Rollback()
			return repository.ErrIDMustNotReferToItself
		}

		var target Department
		if err := tx.First(&target, *reassignTo).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return repository.ErrNotFound
			}
			return err
		}

		inSubtree, err := p.checkParentInSubtreeTx(ctx, tx, id, *reassignTo)
		if err != nil {
			tx.Rollback()
			return err
		}
		if inSubtree {
			tx.Rollback()
			return repository.ErrConflict
		}

		if err := tx.Model(&Employee{}).
			Where("department_id = ?", id).
			Update("department_id", *reassignTo).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(&Department{}).
			Where("parent_id = ?", id).
			Update("parent_id", dep.ParentID).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Delete(&Department{}, id).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit().Error; err != nil {
			return err
		}
		return nil
	}

	if err := tx.Delete(&Department{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
