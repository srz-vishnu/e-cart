package controller

// import (
// 	"e-cart/app/service/mocks"
// 	"e-cart/pkg/e"
// 	"errors"
// 	"fmt"

// 	"net/http/httptest"
// 	"testing"

// 	"github.com/go-playground/assert/v2"
// )

// func TestUpdateUserDetails(t *testing.T) {

// 	userMock := new(mocks.UserService)
// 	con := NewUserController(userMock)

// 	tests := []struct {
// 		name    string
// 		status  int
// 		want    string
// 		Error   error
// 		wantErr bool
// 	}{
// 		{
// 			name:    "success_case",
// 			status:  200,
// 			want:    `{"status":"ok","result":"success"}`,
// 			wantErr: false,
// 		},
// 		{
// 			name: "fail_update",
// 			Error: &e.WrapErrorImpl{
// 				ErrorCode: 400,
// 				Msg:       "Bad Request",
// 				RootCause: errors.New("Invalid Request"),
// 				Loglevel:  "Error",
// 			},
// 			status: 400,
// 			want:   `{"status":"nok","error":{"code":400,"message":"Bad Request","details":["Invalid Request"]}}`,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {

// 			res := httptest.NewRecorder()
// 			req := httptest.NewRequest("Get", "/", nil)
// 			l := log.LoggerFromRequest(req)

// 			// Mocking `UpdateUserDetails()` of Service mock
// 			userMock.Mock.On("UpdateUserDetails", req.Context(), l, req).Once().Return(test.Error)

// 			// calling the function
// 			con.UpdateUserDetails(res, req)

// 			// comparing the response code and body with expected
// 			assert.Equal(t, test.status, res.Code, fmt.Sprintf("Expexted Status code : %d \n Status code we got : %d", test.status, res.Code))
// 			assert.Equal(t, test.want, res.Body.String(), fmt.Sprintf("Expexted Response body : %s \n Response body we got : %s", test.want, res.Body.String()))
// 		})
// 	}
// }
