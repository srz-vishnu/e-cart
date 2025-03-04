package controller

import (
	"e-cart/app/service"
	"e-cart/pkg/api"
	"e-cart/pkg/e"
	"net/http"
)

type UserController interface {
	UserDetails(w http.ResponseWriter, r *http.Request)
	UpdateUserDetails(w http.ResponseWriter, r *http.Request)
}

type UserControllerStructImpl struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerStructImpl{
		userService: userService,
	}
}

func (c *UserControllerStructImpl) UserDetails(w http.ResponseWriter, r *http.Request) {
	resp, err := c.userService.SaveUserDetails(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to create user")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, resp)
}

func (c *UserControllerStructImpl) UpdateUserDetails(w http.ResponseWriter, r *http.Request) {

	err := c.userService.UpdateUserDetails(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to update user")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, "success")
}
