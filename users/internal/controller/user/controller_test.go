package user

import (
	"context"
	itmgate "github.com/JMURv/e-commerce/users/internal/gateway/items"
	repo "github.com/JMURv/e-commerce/users/internal/repository"
	"github.com/JMURv/e-commerce/users/internal/repository/memory"
	"github.com/JMURv/e-commerce/users/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserCreation(t *testing.T) {
	svc := New(memory.New(), itmgate.Gateway{}) // use in-memory db and empty gateway

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
			wantErr:  repo.ErrUsernameIsRequired,
		},
		{
			name:     "noEmail",
			testData: noEmailUserData,
			wantRes:  nil,
			wantErr:  repo.ErrEmailIsRequired,
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
