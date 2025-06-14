// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	dto "e-cart/app/dto"
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// AddItemToCart provides a mock function with given fields: r
func (_m *UserService) AddItemToCart(r *http.Request) (*dto.CartItemResponse, error) {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for AddItemToCart")
	}

	var r0 *dto.CartItemResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request) (*dto.CartItemResponse, error)); ok {
		return rf(r)
	}
	if rf, ok := ret.Get(0).(func(*http.Request) *dto.CartItemResponse); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.CartItemResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddItemsToFavourites provides a mock function with given fields: r
func (_m *UserService) AddItemsToFavourites(r *http.Request) error {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for AddItemsToFavourites")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*http.Request) error); ok {
		r0 = rf(r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ClearCart provides a mock function with given fields: r
func (_m *UserService) ClearCart(r *http.Request) error {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for ClearCart")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*http.Request) error); ok {
		r0 = rf(r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserFavouriteBrands provides a mock function with given fields: r
func (_m *UserService) GetUserFavouriteBrands(r *http.Request) ([]dto.FavoriteBrandResponse, error) {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for GetUserFavouriteBrands")
	}

	var r0 []dto.FavoriteBrandResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request) ([]dto.FavoriteBrandResponse, error)); ok {
		return rf(r)
	}
	if rf, ok := ret.Get(0).(func(*http.Request) []dto.FavoriteBrandResponse); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.FavoriteBrandResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoginUser provides a mock function with given fields: r
func (_m *UserService) LoginUser(r *http.Request) (*dto.LoginResponse, error) {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for LoginUser")
	}

	var r0 *dto.LoginResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request) (*dto.LoginResponse, error)); ok {
		return rf(r)
	}
	if rf, ok := ret.Get(0).(func(*http.Request) *dto.LoginResponse); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.LoginResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrderHistory provides a mock function with given fields: r
func (_m *UserService) OrderHistory(r *http.Request) ([]*dto.ItemOrderedResponse, error) {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for OrderHistory")
	}

	var r0 []*dto.ItemOrderedResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request) ([]*dto.ItemOrderedResponse, error)); ok {
		return rf(r)
	}
	if rf, ok := ret.Get(0).(func(*http.Request) []*dto.ItemOrderedResponse); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*dto.ItemOrderedResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PlaceOrder provides a mock function with given fields: r
func (_m *UserService) PlaceOrder(r *http.Request) (*dto.ItemOrderedResponse, error) {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for PlaceOrder")
	}

	var r0 *dto.ItemOrderedResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request) (*dto.ItemOrderedResponse, error)); ok {
		return rf(r)
	}
	if rf, ok := ret.Get(0).(func(*http.Request) *dto.ItemOrderedResponse); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.ItemOrderedResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveUserDetails provides a mock function with given fields: r
func (_m *UserService) SaveUserDetails(r *http.Request) (*dto.SaveUserResponse, error) {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for SaveUserDetails")
	}

	var r0 *dto.SaveUserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request) (*dto.SaveUserResponse, error)); ok {
		return rf(r)
	}
	if rf, ok := ret.Get(0).(func(*http.Request) *dto.SaveUserResponse); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.SaveUserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUserDetails provides a mock function with given fields: r
func (_m *UserService) UpdateUserDetails(r *http.Request) error {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUserDetails")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*http.Request) error); ok {
		r0 = rf(r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ViewUserCart provides a mock function with given fields: r
func (_m *UserService) ViewUserCart(r *http.Request) ([]*dto.ViewCart, error) {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for ViewUserCart")
	}

	var r0 []*dto.ViewCart
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request) ([]*dto.ViewCart, error)); ok {
		return rf(r)
	}
	if rf, ok := ret.Get(0).(func(*http.Request) []*dto.ViewCart); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*dto.ViewCart)
		}
	}

	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserService(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
