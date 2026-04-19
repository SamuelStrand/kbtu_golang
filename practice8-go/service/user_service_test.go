package service

import (
	"errors"
	"testing"

	"practice8/repository"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().GetByEmail("x@mail.com").Return(&repository.User{ID: 9, Email: "x@mail.com"}, nil)

	err := svc.RegisterUser(&repository.User{ID: 1, Email: "x@mail.com"}, "x@mail.com")
	require.Error(t, err)
	require.Contains(t, err.Error(), "already exists")
}

func TestRegisterUser_NewUserSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	u := &repository.User{ID: 2, Name: "Bob", Email: "bob@mail.com"}

	mockRepo.EXPECT().GetByEmail("bob@mail.com").Return(nil, nil)
	mockRepo.EXPECT().CreateUser(u).Return(nil)

	require.NoError(t, svc.RegisterUser(u, "bob@mail.com"))
}

func TestRegisterUser_RepoErrorOnCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	u := &repository.User{ID: 3, Name: "Eve", Email: "eve@mail.com"}

	mockRepo.EXPECT().GetByEmail("eve@mail.com").Return(nil, nil)
	mockRepo.EXPECT().CreateUser(u).Return(errors.New("db error"))

	err := svc.RegisterUser(u, "eve@mail.com")
	require.Error(t, err)
	require.Contains(t, err.Error(), "db error")
}

func TestUpdateUserName_EmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	err := svc.UpdateUserName(10, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot be empty")
}

func TestUpdateUserName_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().GetUserByID(10).Return(nil, errors.New("not found"))

	err := svc.UpdateUserName(10, "NewName")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

func TestUpdateUserName_Success_VerifyNameChanged(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	u := &repository.User{ID: 10, Name: "OldName", Email: "a@mail.com"}
	mockRepo.EXPECT().GetUserByID(10).Return(u, nil)

	mockRepo.EXPECT().
		UpdateUser(gomock.Any()).
		DoAndReturn(func(updated *repository.User) error {
			require.Equal(t, "NewName", updated.Name) // verify changed
			return nil
		})

	require.NoError(t, svc.UpdateUserName(10, "NewName"))
}

func TestUpdateUserName_UpdateUserFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	u := &repository.User{ID: 10, Name: "OldName", Email: "a@mail.com"}
	mockRepo.EXPECT().GetUserByID(10).Return(u, nil)
	mockRepo.EXPECT().UpdateUser(gomock.Any()).Return(errors.New("update failed"))

	err := svc.UpdateUserName(10, "NewName")
	require.Error(t, err)
	require.Contains(t, err.Error(), "update failed")
}

func TestDeleteUser_AttemptToDeleteAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	err := svc.DeleteUser(1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not allowed")
}

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().DeleteUser(2).Return(nil)
	require.NoError(t, svc.DeleteUser(2))
}

func TestDeleteUser_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().DeleteUser(2).Return(errors.New("repo error"))

	err := svc.DeleteUser(2)
	require.Error(t, err)
	require.Contains(t, err.Error(), "repo error")
}
