package service

// import (
// 	"context"
// 	"e-cart/pkg/e"
// 	"e-cart/pkg/middleware"
// 	"errors"
// 	"fmt"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // MockUserRepo is a mock implementation of internal.UserRepo
// type MockUserRepo struct {
// 	mock.Mock
// }

// func (m *MockUserRepo) IsUserActive(userID int64) (bool, error) {
// 	args := m.Called(userID)
// 	return args.Bool(0), args.Error(1)
// }

// func TestGetUserIDAndCheckStatus(t *testing.T) {
// 	ctx := context.Background()
// 	userID := int64(123)

// 	userMock := new(MockUserRepo)
// 	serv := &userServiceImpl{
// 		userRepo: userMock,
// 	}

// 	tests := []struct {
// 		name       string
// 		userID     int64
// 		isActive   bool
// 		activeErr  error
// 		contextErr error
// 		want       int64
// 		err        error
// 		wantErr    bool
// 	}{
// 		{
// 			name:     "success_case",
// 			userID:   userID,
// 			isActive: true,
// 			want:     userID,
// 			wantErr:  false,
// 		},
// 		{
// 			name:     "user_not_active",
// 			userID:   userID,
// 			isActive: false,
// 			err:      e.NewError(e.ErrUserBlocked, "user is blocked or inactive", nil),
// 			wantErr:  true,
// 		},
// 		{
// 			name:      "error_checking_user_status",
// 			userID:    userID,
// 			activeErr: errors.New("database error"),
// 			err:       e.NewError(e.ErrGetUserDetails, "error while checking user details", errors.New("database error")),
// 			wantErr:   true,
// 		},
// 		{
// 			name:       "error_getting_user_id",
// 			contextErr: errors.New("user ID not found in context"),
// 			err:        e.NewError(e.ErrContextError, "error while getting userId from ctx", errors.New("user ID not found in context")),
// 			wantErr:    true,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			// Create context with or without user ID
// 			testCtx := ctx
// 			if test.contextErr == nil {
// 				testCtx = context.WithValue(ctx, middleware.UserIDKey, test.userID)
// 			}

// 			// Set up mock expectations
// 			if test.contextErr == nil {
// 				userMock.On("IsUserActive", test.userID).Return(test.isActive, test.activeErr).Once()
// 			}

// 			// Call the function
// 			got, err := serv.getUserIDAndCheckStatus(testCtx)

// 			// Assert results
// 			if test.wantErr {
// 				assert.Equal(t, test.err.Error(), err.Error(), test.name+fmt.Sprintf(" FAILED - Mismatch in error: expected %v, got %v", test.err.Error(), err.Error()))
// 			} else {
// 				assert.Equal(t, test.want, got, fmt.Sprintf("Expected result is : %v \n Result we got is : %v", test.want, got))
// 			}

// 			// Verify all expectations were met
// 			userMock.AssertExpectations(t)
// 		})
// 	}
// }
