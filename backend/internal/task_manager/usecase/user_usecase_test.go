package usecase

import (
	"database/sql"
	"errors"
	"github.com/oOSomnus/transflate/internal/task_manager/repository"
	"golang.org/x/crypto/bcrypt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	usecase := NewUserUsecase(mockRepo)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockRepo.EXPECT().FindUsrWithUsername("validUser").Return(string(hashedPassword), nil).AnyTimes()
	mockRepo.EXPECT().FindUsrWithUsername("invalidUser").Return("", sql.ErrNoRows).AnyTimes()

	tests := []struct {
		name     string
		username string
		password string
		want     bool
		wantErr  string
	}{
		{"validCredentials", "validUser", "password123", true, ""},
		{"invalidPassword", "validUser", "wrongPassword", false, ErrInvalidCredentials},
		{"nonExistentUser", "invalidUser", "anyPassword", false, ErrInvalidCredentials},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := usecase.Authenticate(tt.username, tt.password)
				if tt.wantErr != "" {
					assert.EqualError(t, err, tt.wantErr)
				} else {
					assert.NoError(t, err)
				}
				assert.Equal(t, tt.want, got)
			},
		)
	}
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	usecase := NewUserUsecase(mockRepo)

	mockRepo.EXPECT().IfUserExists("existingUser").Return(true, nil).AnyTimes()
	mockRepo.EXPECT().IfUserExists(gomock.Not("existingUser")).Return(false, nil).AnyTimes()
	mockRepo.EXPECT().CreateUser("validUser", "strongPassword1").Return(nil).AnyTimes()
	mockRepo.EXPECT().CreateUser("invalidUser", gomock.Any()).Return(errors.New(ErrCreateUser)).AnyTimes()
	mockRepo.EXPECT().CreateUser("correctPwdUser", gomock.Any()).Return(nil).AnyTimes()
	tests := []struct {
		name     string
		username string
		password string
		wantErr  string
	}{
		{"validUser", "correctPwdUser", "strongPassword1", ""},
		{"existingUser", "existingUser", "strongPassword1", ErrUserAlreadyExists},
		{"shortUsername", "a", "strongPassword1", ErrInvalidRegInfo},
		{"longUsername", "verylongusernamehere", "strongPassword1", ErrInvalidRegInfo},
		{"shortPassword", "validUser", "short", ErrInvalidRegInfo},
		{"emptyFields", "", "", ErrEmptyInput},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				err := usecase.CreateUser(tt.username, tt.password)
				if tt.wantErr != "" {
					assert.EqualError(t, err, tt.wantErr)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}

func TestDecreaseBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	usecase := NewUserUsecase(mockRepo)

	mockRepo.EXPECT().DecreaseBalance("validUser", 50).Return(nil).AnyTimes()
	mockRepo.EXPECT().DecreaseBalance("invalidUser", 50).Return(errors.New("decrease balance error")).AnyTimes()

	tests := []struct {
		name     string
		username string
		balance  int
		wantErr  string
	}{
		{"validDecrease", "validUser", 50, ""},
		{"decreaseError", "invalidUser", 50, "failed to decrease balance: decrease balance error"},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				err := usecase.DecreaseBalance(tt.username, tt.balance)
				if tt.wantErr != "" {
					assert.EqualError(t, err, tt.wantErr)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}

func TestCheckBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	usecase := NewUserUsecase(mockRepo)

	mockRepo.EXPECT().GetBalance("validUser").Return(100, nil).AnyTimes()
	mockRepo.EXPECT().GetBalance("negativeBalanceUser").Return(-10, nil).AnyTimes()
	mockRepo.EXPECT().GetBalance("invalidUser").Return(0, errors.New("balance retrieval error")).AnyTimes()

	tests := []struct {
		name     string
		username string
		want     int
		wantErr  string
	}{
		{"validBalance", "validUser", 100, ""},
		{"negativeBalanceError", "negativeBalanceUser", 0, "negative balance"},
		{"retrievalError", "invalidUser", 0, "failed to get balance: balance retrieval error"},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := usecase.CheckBalance(tt.username)
				if tt.wantErr != "" {
					assert.EqualError(t, err, tt.wantErr)
				} else {
					assert.NoError(t, err)
				}
				assert.Equal(t, tt.want, got)
			},
		)
	}
}
