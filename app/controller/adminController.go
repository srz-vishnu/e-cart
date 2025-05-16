package controller

import (
	"e-cart/app/service"
	"e-cart/pkg/api"
	"e-cart/pkg/e"
	"net/http"
)

type AdminController interface {
	BlockUser(w http.ResponseWriter, r *http.Request)
	UnBlockUser(w http.ResponseWriter, r *http.Request)
	GetAllUserDetail(w http.ResponseWriter, r *http.Request)
	CustomerOrderHistoryById(w http.ResponseWriter, r *http.Request)
	GetAllBlockedUserDetail(w http.ResponseWriter, r *http.Request)
	CustomerOrderHistory(w http.ResponseWriter, r *http.Request)
}

type AdminControlImpl struct {
	adminService service.AdminService
}

func NewAdminController(adminService service.AdminService) AdminController {
	return &AdminControlImpl{
		adminService: adminService,
	}
}

func (c *AdminControlImpl) BlockUser(w http.ResponseWriter, r *http.Request) {

	err := c.adminService.BlockUser(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to block user")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}

	api.Success(w, http.StatusOK, "Successfully blocked user")
}

func (c *AdminControlImpl) UnBlockUser(w http.ResponseWriter, r *http.Request) {

	err := c.adminService.UnBlockUser(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to unblock user")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}

	api.Success(w, http.StatusOK, "Successfully unblocked user")
}

func (c *AdminControlImpl) GetAllUserDetail(w http.ResponseWriter, r *http.Request) {

	resp, err := c.adminService.GetAllUserDetail(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to get all userdetails ")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}

	api.Success(w, http.StatusOK, resp)
}

func (c *AdminControlImpl) GetAllBlockedUserDetail(w http.ResponseWriter, r *http.Request) {

	resp, err := c.adminService.GetAllBlockedUserDetail(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to get blocked userdetails ")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}

	api.Success(w, http.StatusOK, resp)
}

func (c *AdminControlImpl) CustomerOrderHistoryById(w http.ResponseWriter, r *http.Request) {
	resp, err := c.adminService.CustomerOrderHistoryById(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to get the customer order history")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, resp)
}

func (c *AdminControlImpl) CustomerOrderHistory(w http.ResponseWriter, r *http.Request) {
	resp, err := c.adminService.CustomerOrderHistory(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to get the customer order history")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, resp)
}
