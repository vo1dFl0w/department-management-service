package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vo1dFl0w/department-management-service/internal/domain"
	"github.com/vo1dFl0w/department-management-service/internal/repository"
	"github.com/vo1dFl0w/department-management-service/internal/test/mocks"
	"github.com/vo1dFl0w/department-management-service/internal/usecase"
)

var (
	id       int    = 3
	longName string = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip!"
)

func TestDepartmentManagerService_CreateNewDepartment(t *testing.T) {
	testCases := []struct {
		name          string
		nameParam     string
		parentIDParam *int64
		hasRepoErr    bool
		expErrMessage error
		expErr        bool
	}{
		{
			name:          "valid",
			nameParam:     "department",
			parentIDParam: nil,
			hasRepoErr:    false,
			expErrMessage: nil,
			expErr:        false,
		},
		{
			name:          "valid with parent id",
			nameParam:     "department",
			parentIDParam: ptr(int64(id)),
			hasRepoErr:    false,
			expErrMessage: nil,
			expErr:        false,
		},
		{
			name:          "empty name param",
			nameParam:     "",
			parentIDParam: ptr(int64(id)),
			hasRepoErr:    false,
			expErrMessage: domain.ErrEmptyDepartmentName,
			expErr:        true,
		},
		{
			name:          "to long name param",
			nameParam:     longName,
			parentIDParam: ptr(int64(id)),
			hasRepoErr:    false,
			expErrMessage: domain.ErrToLongDepartmentName,
			expErr:        true,
		},
		{
			name:          "conflict",
			nameParam:     "department",
			parentIDParam: ptr(int64(id)),
			hasRepoErr:    true,
			expErrMessage: repository.ErrConflict,
			expErr:        true,
		},
		{
			name:          "not found",
			nameParam:     "department",
			parentIDParam: ptr(int64(id)),
			hasRepoErr:    true,
			expErrMessage: repository.ErrConflict,
			expErr:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			postgresRepo := &mocks.PostgresMock{}
			depSrv := usecase.NewService(postgresRepo)

			if tc.expErr {
				if tc.hasRepoErr {
					postgresRepo.On("CreateDepartment", mock.Anything, tc.nameParam, tc.parentIDParam).Return(nil, tc.expErrMessage).Once()

					_, err := depSrv.CreateNewDepartment(context.Background(), tc.nameParam, tc.parentIDParam)
					assert.Error(t, err)
					assert.ErrorContains(t, err, tc.expErrMessage.Error())

					postgresRepo.AssertExpectations(t)
				} else {
					_, err := depSrv.CreateNewDepartment(context.Background(), tc.nameParam, tc.parentIDParam)
					assert.Error(t, err)
					assert.ErrorContains(t, err, tc.expErrMessage.Error())

					postgresRepo.AssertNotCalled(t, "CreateDepartment", mock.Anything, tc.nameParam, tc.parentIDParam)
				}
			} else {
				dep := &domain.Department{
					ID:        5,
					Name:      tc.nameParam,
					ParentID:  tc.parentIDParam,
					CreatedAt: time.Now().UTC(),
				}

				postgresRepo.On("CreateDepartment", mock.Anything, tc.nameParam, tc.parentIDParam).Return(dep, nil).Once()

				res, err := depSrv.CreateNewDepartment(context.Background(), tc.nameParam, tc.parentIDParam)
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, dep, res)

				postgresRepo.AssertExpectations(t)
			}
		})
	}
}

func TestDepartmentManagerService_CreateNewEmployee(t *testing.T) {
	hiredAt := time.Now().UTC()

	testCases := []struct {
		name          string
		depID         int64
		fullName      string
		position      string
		hiredAt       time.Time
		hasRepoErr    bool
		expErrMessage error
		expErr        bool
	}{
		{
			name:          "valid",
			depID:         int64(id),
			fullName:      "Jown Joe",
			position:      "backend developer",
			hiredAt:       hiredAt,
			hasRepoErr:    false,
			expErrMessage: nil,
			expErr:        false,
		},
		{
			name:          "valid with parent id",
			depID:         int64(id),
			fullName:      "Jown Joe",
			position:      "backend developer",
			hiredAt:       hiredAt,
			hasRepoErr:    false,
			expErrMessage: nil,
			expErr:        false,
		},
		{
			name:          "empty full name",
			depID:         int64(id),
			fullName:      "",
			position:      "backend developer",
			hiredAt:       hiredAt,
			hasRepoErr:    false,
			expErrMessage: domain.ErrEmptyFullName,
			expErr:        true,
		},
		{
			name:          "too long full name",
			depID:         int64(id),
			fullName:      longName,
			position:      "backend developer",
			hiredAt:       hiredAt,
			hasRepoErr:    false,
			expErrMessage: domain.ErrToLongFullName,
			expErr:        true,
		},
		{
			name:          "empty position name",
			depID:         int64(id),
			fullName:      "Jown Doe",
			position:      "",
			hiredAt:       hiredAt,
			hasRepoErr:    false,
			expErrMessage: domain.ErrEmptyPositionName,
			expErr:        true,
		},
		{
			name:          "too long position name",
			depID:         int64(id),
			fullName:      "Jown Doe",
			position:      longName,
			hiredAt:       hiredAt,
			hasRepoErr:    false,
			expErrMessage: domain.ErrToLongPositionName,
			expErr:        true,
		},
		{
			name:          "not found",
			depID:         int64(id),
			fullName:      "Jown Doe",
			position:      "backend developer",
			hiredAt:       hiredAt,
			hasRepoErr:    true,
			expErrMessage: repository.ErrNotFound,
			expErr:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			postgresRepo := &mocks.PostgresMock{}
			depSrv := usecase.NewService(postgresRepo)

			if tc.expErr {
				if !tc.hasRepoErr {
					_, err := depSrv.CreateNewEmployee(context.Background(), tc.depID, tc.fullName, tc.position, &tc.hiredAt)
					assert.Error(t, err)
					assert.ErrorContains(t, err, tc.expErrMessage.Error())

					postgresRepo.AssertNotCalled(t, "CreateEmployee", mock.Anything, tc.depID, tc.fullName, tc.position, &tc.hiredAt)
				} else {
					postgresRepo.On("CreateEmployee", mock.Anything, tc.depID, tc.fullName, tc.position, &tc.hiredAt).Return(nil, tc.expErrMessage).Once()

					_, err := depSrv.CreateNewEmployee(context.Background(), tc.depID, tc.fullName, tc.position, &tc.hiredAt)
					assert.Error(t, err)
					assert.ErrorContains(t, err, tc.expErrMessage.Error())

					postgresRepo.AssertExpectations(t)
				}
			} else {
				emp := &domain.Employee{
					ID:           1,
					DepartmentID: tc.depID,
					FullName:     tc.fullName,
					Position:     tc.position,
					HiredAt:      &tc.hiredAt,
					CreatedAt:    time.Now().UTC(),
				}

				postgresRepo.On("CreateEmployee", mock.Anything, tc.depID, tc.fullName, tc.position, &tc.hiredAt).Return(emp, nil).Once()

				res, err := depSrv.CreateNewEmployee(context.Background(), tc.depID, tc.fullName, tc.position, &tc.hiredAt)
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, emp, res)

				postgresRepo.AssertExpectations(t)
			}
		})
	}
}

func TestDepartmentManagerService_GetDepartmentByID(t *testing.T) {
	testCases := []struct {
		name             string
		id               int64
		depth            int
		actualDepth      int
		includeEmployees bool
		expErrMessage    error
		expErr           bool
	}{
		{
			name:             "success",
			id:               int64(id),
			depth:            1,
			actualDepth:      1,
			includeEmployees: true,
			expErrMessage:    nil,
			expErr:           false,
		},
		{
			name:             "success (depth < 0)",
			id:               int64(id),
			depth:            0,
			actualDepth:      1,
			includeEmployees: true,
			expErrMessage:    nil,
			expErr:           false,
		},
		{
			name:             "success (depth > 5)",
			id:               int64(id),
			depth:            6,
			actualDepth:      5,
			includeEmployees: true,
			expErrMessage:    nil,
			expErr:           false,
		},
		{
			name:             "not found",
			id:               int64(id),
			depth:            5,
			actualDepth:      5,
			includeEmployees: true,
			expErrMessage:    domain.ErrNotFound,
			expErr:           true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			postgresRepo := &mocks.PostgresMock{}
			depSrv := usecase.NewService(postgresRepo)

			if !tc.expErr {
				dep := &domain.DepartmentTree{
					Department: domain.Department{
						ID:        5,
						Name:      "Backend department",
						ParentID:  ptr(int64(id)),
						CreatedAt: time.Now().UTC(),
					},
					Employees: []domain.Employee{{
						ID:           1,
						DepartmentID: int64(id),
						FullName:     "Jown Doe",
						Position:     "backend developer",
						HiredAt:      nil,
						CreatedAt:    time.Now().UTC(),
					}},
					Children: []*domain.DepartmentTree{},
				}

				postgresRepo.On("GetDepartment", mock.Anything, tc.id, tc.actualDepth, tc.includeEmployees).Return(dep, nil).Once()

				res, err := depSrv.GetDepartmentByID(context.Background(), tc.id, tc.depth, tc.includeEmployees)
				assert.NoError(t, err)
				assert.NotNil(t, res)

				postgresRepo.AssertExpectations(t)
			} else {
				postgresRepo.On("GetDepartment", mock.Anything, tc.id, tc.actualDepth, tc.includeEmployees).Return(nil, repository.ErrNotFound).Once()

				_, err := depSrv.GetDepartmentByID(context.Background(), tc.id, tc.depth, tc.includeEmployees)
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expErrMessage.Error())

				postgresRepo.AssertExpectations(t)
			}
		})
	}
}

func TestDepartmentManagerService_UpdateDepartmentByID(t *testing.T) {
	testCases := []struct {
		name          string
		id            int64
		namePtr       *string
		parentID      *int64
		hasRepoErr    bool
		expErrMessage error
		expErr        bool
	}{
		{
			name:          "success (all params)",
			id:            int64(id),
			namePtr:       ptr("Backend department"),
			parentID:      ptr(int64(1)),
			hasRepoErr:    false,
			expErrMessage: nil,
			expErr:        false,
		},
		{
			name:          "success (with namePtr)",
			id:            int64(id),
			namePtr:       ptr("Backend department"),
			parentID:      nil,
			hasRepoErr:    false,
			expErrMessage: nil,
			expErr:        false,
		},
		{
			name:          "id must not refer to it self",
			id:            int64(id),
			namePtr:       ptr("Backend department"),
			parentID:      ptr(int64(id)),
			hasRepoErr:    false,
			expErrMessage: domain.ErrIDMustNotReferToItself,
			expErr:        true,
		},
		{
			name:          "too long department name",
			id:            int64(id),
			namePtr:       ptr(longName),
			parentID:      ptr(int64(1)),
			hasRepoErr:    false,
			expErrMessage: domain.ErrToLongDepartmentName,
			expErr:        true,
		},
		{
			name:          "conflict",
			id:            int64(id),
			namePtr:       ptr("Backend department"),
			parentID:      ptr(int64(1)),
			hasRepoErr:    true,
			expErrMessage: repository.ErrConflict,
			expErr:        true,
		},
		{
			name:          "not found",
			id:            int64(id),
			namePtr:       ptr("Backend department"),
			parentID:      ptr(int64(1)),
			hasRepoErr:    true,
			expErrMessage: repository.ErrNotFound,
			expErr:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			postgresRepo := &mocks.PostgresMock{}
			depSrv := usecase.NewService(postgresRepo)

			if tc.expErr {
				if !tc.hasRepoErr {
					_, err := depSrv.UpdateDepartmentByID(context.Background(), tc.id, tc.namePtr, tc.parentID)
					assert.Error(t, err)
					assert.ErrorContains(t, err, tc.expErrMessage.Error())

					postgresRepo.AssertNotCalled(t, "UpdateDepartment", mock.Anything, tc.id, tc.namePtr, tc.parentID)
				} else {
					postgresRepo.On("UpdateDepartment", mock.Anything, tc.id, tc.namePtr, tc.parentID).Return(nil, tc.expErrMessage).Once()

					_, err := depSrv.UpdateDepartmentByID(context.Background(), tc.id, tc.namePtr, tc.parentID)
					assert.Error(t, err)
					assert.ErrorContains(t, err, tc.expErrMessage.Error())

					postgresRepo.AssertExpectations(t)
				}
			} else {
				dep := &domain.Department{
					ID:        tc.id,
					Name:      "Frontend department",
					ParentID:  tc.parentID,
					CreatedAt: time.Now().UTC(),
				}

				postgresRepo.On("UpdateDepartment", mock.Anything, tc.id, tc.namePtr, tc.parentID).Return(dep, nil).Once()

				res, err := depSrv.UpdateDepartmentByID(context.Background(), tc.id, tc.namePtr, tc.parentID)
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, dep, res)

				postgresRepo.AssertExpectations(t)
			}
		})
	}
}

func TestDepartmentManagerService_DeleteDepartmentByID(t *testing.T) {
	testCases := []struct {
		name          string
		id            int64
		mode          string
		reassignTo    *int64
		hasRepoErr    bool
		expErrMessage error
		expErr        bool
	}{
		{
			name:          "success (cascade)",
			id:            int64(id),
			mode:          "cascade",
			reassignTo:    nil,
			hasRepoErr:    false,
			expErrMessage: nil,
			expErr:        false,
		},
		{
			name:          "success (reassign)",
			id:            int64(id),
			mode:          "reassign",
			reassignTo:    ptr(int64(1)),
			hasRepoErr:    false,
			expErrMessage: nil,
			expErr:        false,
		},
		{
			name:          "reassign id must be set",
			id:            int64(id),
			mode:          "reassign",
			reassignTo:    nil,
			hasRepoErr:    false,
			expErrMessage: domain.ErrResignIDMustBeSet,
			expErr:        true,
		},
		{
			name:          "id must not refer to itself",
			id:            int64(id),
			mode:          "reassign",
			reassignTo:    ptr(int64(id)),
			hasRepoErr:    false,
			expErrMessage: domain.ErrIDMustNotReferToItself,
			expErr:        true,
		},
		{
			name:          "conflict",
			id:            int64(id),
			mode:          "reassign",
			reassignTo:    ptr(int64(1)),
			hasRepoErr:    true,
			expErrMessage: repository.ErrConflict,
			expErr:        true,
		},
		{
			name:          "not found",
			id:            int64(id),
			mode:          "reassign",
			reassignTo:    ptr(int64(1)),
			hasRepoErr:    true,
			expErrMessage: repository.ErrNotFound,
			expErr:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			postgresRepo := &mocks.PostgresMock{}
			depSrv := usecase.NewService(postgresRepo)

			if tc.expErr {
				if !tc.hasRepoErr {
					err := depSrv.DeleteDepartmentByID(context.Background(), tc.id, tc.mode, tc.reassignTo)
					assert.Error(t, err)
					assert.ErrorContains(t, err, tc.expErrMessage.Error())

					postgresRepo.AssertNotCalled(t, "DeleteDepartment", mock.Anything, tc.id, tc.mode, tc.reassignTo)
				} else {
					postgresRepo.On("DeleteDepartment", mock.Anything, tc.id, tc.mode, tc.reassignTo).Return(tc.expErrMessage).Once()

					err := depSrv.DeleteDepartmentByID(context.Background(), tc.id, tc.mode, tc.reassignTo)
					assert.Error(t, err)
					assert.ErrorContains(t, err, tc.expErrMessage.Error())

					postgresRepo.AssertExpectations(t)
				}
			} else {
				postgresRepo.On("DeleteDepartment", mock.Anything, tc.id, tc.mode, tc.reassignTo).Return(nil).Once()

				err := depSrv.DeleteDepartmentByID(context.Background(), tc.id, tc.mode, tc.reassignTo)
				assert.NoError(t, err)

				postgresRepo.AssertExpectations(t)
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
