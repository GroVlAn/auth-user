package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/GroVlAn/auth-user/internal/core"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Create(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	validUser := core.User{
		Username: "john_doe",
		Email:    "example@example.com",
		Password: "12345WWw##3",
		Fullname: "John Doe",
	}

	tests := []struct {
		name      string
		user      core.User
		setupMock func(m *mockrepo, h *mockhasher)
		check     func(t *testing.T, err error, m *mockrepo, h *mockhasher)
	}{
		{
			name: "validation empty user error",
			user: core.User{},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
				m.AssertNotCalled(t, "Exist", mock.Anything, mock.Anything)
				h.AssertNotCalled(t, "Hash", mock.Anything, mock.Anything)
			},
		},
		{
			name: "validation empty username error",
			user: core.User{
				Email:    "example@example.com",
				Password: "12345WWw##3",
				Fullname: "John Doe",
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
				m.AssertNotCalled(t, "Exist", mock.Anything, mock.Anything)
				h.AssertNotCalled(t, "Hash", mock.Anything, mock.Anything)
			},
		},
		{
			name: "validation empty email error",
			user: core.User{
				Email:    "example@example.com",
				Password: "12345WWw##3",
				Fullname: "John Doe",
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
				m.AssertNotCalled(t, "Exist", mock.Anything, mock.Anything)
				h.AssertNotCalled(t, "Hash", mock.Anything, mock.Anything)
			},
		},
		{
			name: "validation empty password error",
			user: core.User{
				Username: "john_doe",
				Email:    "example@example.com",
				Fullname: "John Doe",
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
				m.AssertNotCalled(t, "Exist", mock.Anything, mock.Anything)
				h.AssertNotCalled(t, "Hash", mock.Anything, mock.Anything)
			},
		},
		{
			name: "validation empty fullname error",
			user: core.User{
				Username: "john_doe",
				Email:    "example@example.com",
				Password: "12345WWw##3",
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
				m.AssertNotCalled(t, "Exist", mock.Anything, mock.Anything)
				h.AssertNotCalled(t, "Hash", mock.Anything, mock.Anything)
			},
		},
		{
			name: "validation bad email error",
			user: core.User{
				Username: "john_doe",
				Email:    "example",
				Password: "12345WWw##3",
				Fullname: "John Doe",
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
				m.AssertNotCalled(t, "Exist", mock.Anything, mock.Anything)
				h.AssertNotCalled(t, "Hash", mock.Anything, mock.Anything)
			},
		},
		{
			name: "validation bad password error",
			user: core.User{
				Username: "john_doe",
				Email:    "examplee@example.com",
				Password: "12345",
				Fullname: "John Doe",
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
				m.AssertNotCalled(t, "Exist", mock.Anything, mock.Anything)
				h.AssertNotCalled(t, "Hash", mock.Anything, mock.Anything)
			},
		},
		{
			name: "validation bad fullname error",
			user: core.User{
				Username: "john_doe",
				Email:    "examplee@example.com",
				Password: "12345WWw##3",
				Fullname: "John",
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
				m.AssertNotCalled(t, "Exist", mock.Anything, mock.Anything)
				h.AssertNotCalled(t, "Hash", mock.Anything, mock.Anything)
			},
		},
		{
			name: "exist returns error",
			user: validUser,
			setupMock: func(m *mockrepo, h *mockhasher) {
				m.On("Exist", mock.Anything, mock.Anything).
					Return(false, fmt.Errorf("db error")).Once()
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
				m.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
				h.AssertNotCalled(t, "Hash", mock.Anything, mock.Anything)
			},
		},
		{
			name: "user already exists",
			user: validUser,
			setupMock: func(m *mockrepo, h *mockhasher) {
				m.On("Exist", mock.Anything, mock.Anything).
					Return(true, nil).Once()
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
				m.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
			},
		},
		{
			name: "create fails",
			user: validUser,
			setupMock: func(m *mockrepo, h *mockhasher) {
				m.On("Exist", mock.Anything, mock.Anything).
					Return(false, nil).Once()

				h.On("Hash", validUser.Password).
					Return("hashed_password", nil).Once()

				m.On("Create", mock.Anything, mock.Anything).
					Return(fmt.Errorf("db error")).Once()
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
			},
		},
		{
			name: "success",
			user: validUser,
			setupMock: func(m *mockrepo, h *mockhasher) {
				m.On("Exist", mock.Anything, mock.Anything).
					Return(false, nil).Once()

				h.On("Hash", validUser.Password).
					Return("hashed_password", nil).Once()

				m.On("Create", mock.Anything, mock.MatchedBy(func(u core.User) bool {
					return u.ID != "" &&
						u.PasswordHash != "" &&
						u.PasswordHash != u.Password &&
						!u.CreatedAt.IsZero()
				})).Return(nil).Once()
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo)
			hasherRepo := new(mockhasher)
			s := New(mockRepo, hasherRepo)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo, hasherRepo)
			}

			err := s.Create(ctx, tt.user)

			tt.check(t, err, mockRepo, hasherRepo)

			hasherRepo.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_User(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	validQuery := core.UserQuery{
		Username: "john_doe",
	}

	expectedUser := core.User{
		ID:       "123",
		Username: "john_doe",
		Email:    "example@example.com",
		Fullname: "John Doe",
	}

	tests := []struct {
		name      string
		query     core.UserQuery
		setupMock func(m *mockrepo)
		check     func(t *testing.T, user core.User, err error, m *mockrepo)
	}{
		{
			name:  "validation error - empty query",
			query: core.UserQuery{},
			check: func(t *testing.T, user core.User, err error, m *mockrepo) {
				require.Error(t, err)
				require.Equal(t, core.User{}, user)
				m.AssertNotCalled(t, mock.Anything, mock.Anything)
			},
		},
		{
			name:  "repository returns error",
			query: validQuery,
			setupMock: func(m *mockrepo) {
				m.On("User", mock.Anything, validQuery).
					Return(core.User{}, fmt.Errorf("db error")).Once()
			},
			check: func(t *testing.T, user core.User, err error, m *mockrepo) {
				require.Error(t, err)
				require.Equal(t, core.User{}, user)
			},
		},
		{
			name:  "success",
			query: validQuery,
			setupMock: func(m *mockrepo) {
				m.On("User", mock.Anything, validQuery).
					Return(expectedUser, nil).Once()
			},
			check: func(t *testing.T, user core.User, err error, m *mockrepo) {
				require.NoError(t, err)
				require.Equal(t, expectedUser, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo)
			hasherRepo := new(mockhasher)
			s := New(mockRepo, hasherRepo)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			user, err := s.User(ctx, tt.query)

			tt.check(t, user, err, mockRepo)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_UserInfo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	validQuery := core.UserQuery{
		Username: "john_doe",
	}

	expectedUser := core.UserInfo{
		Username: "john_doe",
		Email:    "example@example.com",
		Fullname: "John Doe",
	}

	tests := []struct {
		name      string
		query     core.UserQuery
		setupMock func(m *mockrepo)
		check     func(t *testing.T, user core.UserInfo, err error, m *mockrepo)
	}{
		{
			name:  "validation error - empty query",
			query: core.UserQuery{},
			check: func(t *testing.T, userInfo core.UserInfo, err error, m *mockrepo) {
				require.Error(t, err)
				require.Equal(t, core.UserInfo{}, userInfo)
				m.AssertNotCalled(t, mock.Anything, mock.Anything)
			},
		},
		{
			name:  "repository returns error",
			query: validQuery,
			setupMock: func(m *mockrepo) {
				m.On("UserInfo", mock.Anything, validQuery).
					Return(core.UserInfo{}, fmt.Errorf("db error")).Once()
			},
			check: func(t *testing.T, user core.UserInfo, err error, m *mockrepo) {
				require.Error(t, err)
				require.Equal(t, core.UserInfo{}, user)
			},
		},
		{
			name:  "success",
			query: validQuery,
			setupMock: func(m *mockrepo) {
				m.On("UserInfo", mock.Anything, validQuery).
					Return(expectedUser, nil).Once()
			},
			check: func(t *testing.T, user core.UserInfo, err error, m *mockrepo) {
				require.NoError(t, err)
				require.Equal(t, expectedUser, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo)
			hasherRepo := new(mockhasher)
			s := New(mockRepo, hasherRepo)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			userInfo, err := s.UserInfo(ctx, tt.query)

			tt.check(t, userInfo, err, mockRepo)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_UpdatePassword(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	validQueryNewPassword := core.UserQueryNewPassword{
		UserQuery: core.UserQuery{
			Username: "john_doe",
		},
		NewPassword: "NewPassword123!",
		OldPassword: "OldPassword123!",
	}

	existingUser := core.User{
		ID:           "123",
		Username:     "john_doe",
		Email:        "example@example.com",
		Fullname:     "John Doe",
		PasswordHash: "OldPassword123!",
	}

	tests := []struct {
		name      string
		query     core.UserQueryNewPassword
		setupMock func(m *mockrepo, h *mockhasher)
		check     func(t *testing.T, err error, m *mockrepo, h *mockhasher)
	}{
		{
			name:  "validation error - empty query",
			query: core.UserQueryNewPassword{},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
			},
		},
		{
			name:  "user not found repository error",
			query: validQueryNewPassword,
			setupMock: func(m *mockrepo, h *mockhasher) {
				m.On("User", mock.Anything, validQueryNewPassword.UserQuery).
					Return(core.User{}, fmt.Errorf("db error")).Once()
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
			},
		},
		{
			name:  "verify old password fails",
			query: validQueryNewPassword,
			setupMock: func(m *mockrepo, h *mockhasher) {
				m.On("User", mock.Anything, validQueryNewPassword.UserQuery).
					Return(existingUser, nil).Once()

				h.On("Compare", existingUser.PasswordHash, validQueryNewPassword.OldPassword).
					Return(fmt.Errorf("wrong password")).Once()
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
			},
		},
		{
			name: "validate new password fails",
			query: core.UserQueryNewPassword{
				UserQuery: core.UserQuery{
					Username: "john_doe",
				},
				NewPassword: "short",
				OldPassword: "OldPassword123!",
			},
			setupMock: func(m *mockrepo, h *mockhasher) {
				m.On("User", mock.Anything, mock.Anything).
					Return(existingUser, nil).Once()

				h.On("Compare", existingUser.PasswordHash, validQueryNewPassword.OldPassword).
					Return(nil).Once()
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
			},
		},
		{
			name: "verify new password fails - same as old",
			query: core.UserQueryNewPassword{
				UserQuery: core.UserQuery{
					Username: "john_doe",
				},
				NewPassword: "OldPassword123!",
				OldPassword: "OldPassword123!",
			},
			setupMock: func(m *mockrepo, h *mockhasher) {
				m.On("User", mock.Anything, mock.Anything).
					Return(existingUser, nil).Once()

				h.On("Compare", existingUser.PasswordHash, validQueryNewPassword.OldPassword).
					Return(nil).Twice()
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
			},
		},
		{
			name:  "repository UpdatePassword returns error",
			query: validQueryNewPassword,
			setupMock: func(m *mockrepo, h *mockhasher) {
				m.On("User", mock.Anything, validQueryNewPassword.UserQuery).
					Return(existingUser, nil).Once()

				h.On("Compare", existingUser.PasswordHash, validQueryNewPassword.OldPassword).
					Return(nil).Once()
				h.On("Compare", validQueryNewPassword.OldPassword, validQueryNewPassword.NewPassword).
					Return(fmt.Errorf("error not same")).Once()

				h.On("Hash", validQueryNewPassword.NewPassword).
					Return("NewPasswordHash", nil).Once()

				m.On("UpdatePassword", mock.Anything, existingUser.ID, "NewPasswordHash").
					Return(fmt.Errorf("db error")).Once()
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.Error(t, err)
			},
		},
		{
			name:  "success",
			query: validQueryNewPassword,
			setupMock: func(m *mockrepo, h *mockhasher) {
				m.On("User", mock.Anything, validQueryNewPassword.UserQuery).
					Return(existingUser, nil).Once()

				h.On("Compare", existingUser.PasswordHash, validQueryNewPassword.OldPassword).
					Return(nil).Once()
				h.On("Compare", validQueryNewPassword.OldPassword, validQueryNewPassword.NewPassword).
					Return(fmt.Errorf("error not same")).Once()

				h.On("Hash", validQueryNewPassword.NewPassword).
					Return("NewPasswordHash", nil).Once()

				m.On("UpdatePassword", mock.Anything, existingUser.ID, "NewPasswordHash").
					Return(nil).Once()
			},
			check: func(t *testing.T, err error, m *mockrepo, h *mockhasher) {
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo)
			hasherRepo := new(mockhasher)
			s := New(mockRepo, hasherRepo)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo, hasherRepo)
			}

			err := s.UpdatePassword(ctx, tt.query)

			tt.check(t, err, mockRepo, hasherRepo)

			hasherRepo.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}
