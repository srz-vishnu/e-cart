package controller

import (
	"e-cart/app/service"
	"e-cart/pkg/api"
	"e-cart/pkg/e"
	"net/http"
)

type UserController interface {
	LoginUser(w http.ResponseWriter, r *http.Request)
	UserDetails(w http.ResponseWriter, r *http.Request)
	UpdateUserDetails(w http.ResponseWriter, r *http.Request)
	ViewUserCart(w http.ResponseWriter, r *http.Request)
	ClearCart(w http.ResponseWriter, r *http.Request)
	AddItemsToCart(w http.ResponseWriter, r *http.Request)
	PlaceOrder(w http.ResponseWriter, r *http.Request)
	OrderHistory(w http.ResponseWriter, r *http.Request)
	AddItemsToFavourites(w http.ResponseWriter, r *http.Request)
	GetUserFavouriteItems(w http.ResponseWriter, r *http.Request)
}

type UserControllerImpl struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{
		userService: userService,
	}
}

func (c *UserControllerImpl) UserDetails(w http.ResponseWriter, r *http.Request) {
	resp, err := c.userService.SaveUserDetails(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to create user")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, resp)
}

func (c *UserControllerImpl) LoginUser(w http.ResponseWriter, r *http.Request) {
	resp, err := c.userService.LoginUser(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to login user")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, resp)
}

func (c *UserControllerImpl) UpdateUserDetails(w http.ResponseWriter, r *http.Request) {

	err := c.userService.UpdateUserDetails(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to update user")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, "success")
}

func (c *UserControllerImpl) AddItemsToCart(w http.ResponseWriter, r *http.Request) {
	resp, err := c.userService.AddItemToCart(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to add items to the cart")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, resp)
}

func (c *UserControllerImpl) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	resp, err := c.userService.PlaceOrder(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to place the order")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, resp)
}

func (c *UserControllerImpl) ViewUserCart(w http.ResponseWriter, r *http.Request) {
	resp, err := c.userService.ViewUserCart(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to view the cart")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, resp)
}

func (c *UserControllerImpl) ClearCart(w http.ResponseWriter, r *http.Request) {
	err := c.userService.ClearCart(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to clear the cart")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, "success")
}

func (c *UserControllerImpl) OrderHistory(w http.ResponseWriter, r *http.Request) {
	resp, err := c.userService.OrderHistory(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to get the order history")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, resp)
}

func (c *UserControllerImpl) AddItemsToFavourites(w http.ResponseWriter, r *http.Request) {
	err := c.userService.AddItemsToFavourites(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to update brand in to favourites")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, "success")
}

func (c *UserControllerImpl) GetUserFavouriteItems(w http.ResponseWriter, r *http.Request) {
	brands, err := c.userService.GetUserFavouriteBrands(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to fetch favourite brands of user")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}

	api.Success(w, http.StatusOK, brands)
}
