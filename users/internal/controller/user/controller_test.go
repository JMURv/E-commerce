package user

import (
	"context"
	registry "github.com/JMURv/e-commerce/pkg/discovery/memory"
	broker "github.com/JMURv/e-commerce/users/internal/broker/memory"
	cache "github.com/JMURv/e-commerce/users/internal/cache/memory"
	itmgate "github.com/JMURv/e-commerce/users/internal/gateway/items"
	repoerrs "github.com/JMURv/e-commerce/users/internal/repository"
	repo "github.com/JMURv/e-commerce/users/internal/repository/memory"
	"github.com/JMURv/e-commerce/users/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var svc *Controller

func init() {
	r := repo.New()
	c := cache.New()
	b := broker.New()
	itemGate := itmgate.New(registry.NewRegistry())
	svc = New(r, c, b, itemGate)
}

func TestUserCreation(t *testing.T) {
	successUserData := &model.User{
		Username: "Test Username",
		Email:    "test@email.com",
	}
	emptyUserData := &model.User{
		Username: "",
		Email:    "",
	}
	noEmailUserData := &model.User{
		Username: "Test Username",
		Email:    "",
	}

	tests := []struct {
		name       string
		testData   *model.User
		expRepoErr error
		wantRes    *model.User
		wantErr    error
	}{
		{
			name:     "success",
			testData: successUserData,
			wantRes:  successUserData,
			wantErr:  nil,
		},
		{
			name:     "empty",
			testData: emptyUserData,
			wantRes:  nil,
			wantErr:  repoerrs.ErrUsernameIsRequired,
		},
		{
			name:     "noEmail",
			testData: noEmailUserData,
			wantRes:  nil,
			wantErr:  repoerrs.ErrEmailIsRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := svc.CreateUser(context.Background(), tt.testData)
			assert.Equal(t, tt.wantRes, u, tt.name)
			assert.Equal(t, tt.wantErr, err, tt.name)
		})
	}
}

func TestUserUpdate(t *testing.T) {
	u, _ := svc.CreateUser(context.Background(), &model.User{
		Username: "Test Username",
		Email:    "test@email.com",
	})
	usrID := u.ID

	successUserData := &model.User{
		Username: "Updated Username",
		Email:    "updated@email.com",
	}
	emptyUserData := &model.User{
		Username: "",
		Email:    "",
	}
	noEmailUserData := &model.User{
		Username: "Test Username",
		Email:    "",
	}

	tests := []struct {
		name       string
		testData   *model.User
		expRepoErr error
		wantRes    *model.User
		wantErr    error
	}{
		{
			name:     "success",
			testData: successUserData,
			wantRes:  successUserData,
			wantErr:  nil,
		},
		{
			name:     "empty",
			testData: emptyUserData,
			wantRes:  nil,
			wantErr:  repoerrs.ErrUsernameIsRequired,
		},
		{
			name:     "noEmail",
			testData: noEmailUserData,
			wantRes:  nil,
			wantErr:  repoerrs.ErrEmailIsRequired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := svc.UpdateUser(context.Background(), usrID, tt.testData)
			assert.Equal(t, tt.wantRes, u, tt.name)
			assert.Equal(t, tt.wantErr, err, tt.name)
		})
	}
}

func TestUserDeletion(t *testing.T) {
	usr, _ := svc.repo.Create(context.Background(), &model.User{
		Username: "Existing Username",
		Email:    "existing@email.com",
	})
	usrID := usr.ID

	tests := []struct {
		name       string
		usrID      uint64
		expRepoErr error
		wantRes    *model.User
		wantErr    error
	}{
		{
			name:    "success",
			usrID:   usrID,
			wantRes: nil,
			wantErr: nil,
		},
		{
			name:    "nonExistentUser",
			usrID:   uint64(123),
			wantRes: nil,
			wantErr: repoerrs.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.DeleteUser(context.Background(), tt.usrID)
			assert.Equal(t, tt.wantErr, err, tt.name)

			_, err = svc.GetUserByID(context.Background(), tt.usrID)
			assert.Equal(t, repoerrs.ErrNotFound, err, tt.name)
		})
	}
}

func TestGetUsersList(t *testing.T) {
	tests := []struct {
		name    string
		wantRes int
		wantErr error
	}{
		{
			name:    "success",
			wantRes: 2,
			wantErr: nil,
		},
	}

	testUsers := []*model.User{
		{Username: "User1", Email: "user1@example.com"},
		{Username: "User2", Email: "user2@example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, user := range testUsers {
				_, _ = svc.repo.Create(context.Background(), user)
				time.Sleep(time.Second * 1)
			}

			usersList, err := svc.GetUsersList(context.Background())
			assert.NoError(t, err, "unexpected error")
			assert.NotNil(t, usersList, "users list is nil")
			assert.Equal(t, len(testUsers), len(usersList), "unexpected number of users")
		})
	}

}

func TestGetUserByID(t *testing.T) {
	usr, _ := svc.CreateUser(context.Background(), &model.User{
		Username: "Test Username",
		Email:    "test@email.com",
	})

	u, err := svc.GetUserByID(context.Background(), usr.ID)

	assert.Equal(t, u.Username, usr.Username, "check username")
	assert.NotNil(t, u, "user is nil")
	assert.NoError(t, err, "unexpected error")
}

func TestGetUserByEmail(t *testing.T) {
	usr, _ := svc.CreateUser(context.Background(), &model.User{
		Username: "Test Username",
		Email:    "test@email.com",
	})

	u, err := svc.GetUserByEmail(context.Background(), "test@email.com")

	assert.Equal(t, u.Username, usr.Username, "check username")
	assert.NotNil(t, u, "user is nil")
	assert.NoError(t, err, "unexpected error")
}
