package service

import (
	"e-cart/app/dto"
	"e-cart/app/internal"
	"e-cart/pkg/e"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type AdminService interface {
	BlockUser(r *http.Request) error
	UnBlockUser(r *http.Request) error
	GetAllUserDetail(r *http.Request) ([]*dto.AllUserDetails, error)
	CustomerOrderHistoryById(r *http.Request) ([]*dto.ItemOrderedResponse, error)
	CustomerOrderHistory(r *http.Request) ([]*dto.ItemOrderedResponse, error)
	GetAllBlockedUserDetail(r *http.Request) ([]*dto.AllUserDetails, error)
}

type AdminServiceImpl struct {
	adminRepo internal.AdminRepo
	userRepo  internal.UserRepo
}

func NewAdminService(adminRepo internal.AdminRepo, userRepo internal.UserRepo) AdminService {
	return &AdminServiceImpl{
		adminRepo: adminRepo,
		userRepo:  userRepo,
	}
}

func (s *AdminServiceImpl) BlockUser(r *http.Request) error {
	args := &dto.BlockUserRequest{}

	//parsing
	err := args.Parse(r)
	if err != nil {
		return e.NewError(e.ErrDecodeRequestBody, "error while parsing", err)
	}

	err = s.adminRepo.BlockUser(args.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewError(e.ErrUserNotFound, "user not found in the table", err)
		}
		return e.NewError(e.ErrBlockUser, "failed to block user", err)
	}
	log.Info().Msg("user with ID %d has been successfully blocked")

	return nil
}

func (s *AdminServiceImpl) UnBlockUser(r *http.Request) error {
	args := &dto.BlockUserRequest{}

	//parsing
	err := args.Parse(r)
	if err != nil {
		return e.NewError(e.ErrDecodeRequestBody, "error while parsing", err)
	}

	err = s.adminRepo.UnBlockUser(args.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewError(e.ErrUserNotFound, "user not found in the table", err)
		}
		return e.NewError(e.ErrUnblockUser, "failed to unblock user", err)
	}
	log.Info().Msg("user with ID %d has been successfully unblocked")

	return nil
}

func (s *AdminServiceImpl) GetAllUserDetail(r *http.Request) ([]*dto.AllUserDetails, error) {
	allUserDetails, err := s.adminRepo.GetAllUsers()
	if err != nil {
		return nil, e.NewError(e.ErrGetUserDetails, "failed to get all user details", err)
	}

	var userDetails []*dto.AllUserDetails

	for _, user := range allUserDetails {
		userDetail := &dto.AllUserDetails{
			UserName: user.Username,
			Mail:     user.Mail,
			Address:  user.Address,
			Phone:    user.Phonenumber,
		}

		userDetails = append(userDetails, userDetail)
		for _, details := range userDetails {
			fmt.Printf("User: %+v\n", *details)
		}
	}

	return userDetails, nil
}

func (s *AdminServiceImpl) GetAllBlockedUserDetail(r *http.Request) ([]*dto.AllUserDetails, error) {
	allUserDetails, err := s.adminRepo.GetAllBlockedUsers()
	if err != nil {
		return nil, e.NewError(e.ErrGetUserDetails, "failed to get blocked user details", err)
	}

	var userDetails []*dto.AllUserDetails

	for _, user := range allUserDetails {
		userDetail := &dto.AllUserDetails{
			UserName: user.Username,
			Mail:     user.Mail,
			Address:  user.Address,
			Phone:    user.Phonenumber,
		}

		userDetails = append(userDetails, userDetail)
	}

	return userDetails, nil
}

func (s *AdminServiceImpl) CustomerOrderHistoryById(r *http.Request) ([]*dto.ItemOrderedResponse, error) {
	args := &dto.SearchByCustomerIdRequest{}

	//parsing
	err := args.Parse(r)
	if err != nil {
		return nil, e.NewError(e.ErrDecodeRequestBody, "error while parsing", err)
	}

	isActive, err := s.userRepo.IsUserActive(args.UserId)
	if err != nil {
		return nil, e.NewError(e.ErrGetUserDetails, "failed to check user status", err)
	}

	if !isActive {
		log.Info().Msg("User is not active.")
		return nil, e.NewError(e.ErrUserNotFound, "user is blocked or inactive", nil)
	}
	log.Info().Msg("User is active")

	userDetails, err := s.userRepo.GetUserByID(args.UserId)
	if err != nil {
		return nil, e.NewError(e.ErrGetUserDetails, "failed to get user details", err)
	}

	orderHistory, err := s.userRepo.GetOrderHistoryByUserID(args.UserId)
	if err != nil {
		return nil, e.NewError(e.ErrGetOrderHistory, "failed to get order history", err)
	}

	var responses []*dto.ItemOrderedResponse

	for _, order := range orderHistory {
		var orderItems []dto.OrderItemResponse

		for _, item := range order.Items {
			orderItems = append(orderItems, dto.OrderItemResponse{
				ProductID:  item.ProductID,
				Quantity:   item.Quantity,
				CategoryID: item.Product.CategoryID,
				BrandName:  item.Product.BrandName,
				Price:      item.Price,
			})
		}

		response := &dto.ItemOrderedResponse{
			Items:      orderItems,
			OrderID:    order.ID,
			TotalPrice: order.Total,
			UserDetails: dto.UserDetailsResponse{
				Username:    userDetails.Username,
				Email:       userDetails.Mail,
				PhoneNumber: userDetails.Phonenumber,
				Address:     userDetails.Address,
				Pincode:     userDetails.Pincode,
			},
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *AdminServiceImpl) CustomerOrderHistory(r *http.Request) ([]*dto.ItemOrderedResponse, error) {
	// Get all orders from the repository
	orders, err := s.adminRepo.GetAllOrders()
	if err != nil {
		return nil, e.NewError(e.ErrGetOrderHistory, "failed to get all orders", err)
	}

	var responses []*dto.ItemOrderedResponse

	// Process each order
	for _, order := range orders {
		var orderItems []dto.OrderItemResponse

		// Process each item in the order
		for _, item := range order.Items {
			orderItems = append(orderItems, dto.OrderItemResponse{
				ProductID:  item.ProductID,
				Quantity:   item.Quantity,
				CategoryID: item.Product.CategoryID,
				BrandName:  item.Product.BrandName,
				Price:      item.Price,
			})
		}

		// Create response for this order
		response := &dto.ItemOrderedResponse{
			Items:      orderItems,
			OrderID:    order.ID,
			TotalPrice: order.Total,
			UserDetails: dto.UserDetailsResponse{
				Username:    order.User.Username,
				Email:       order.User.Mail,
				PhoneNumber: order.User.Phonenumber,
				Address:     order.User.Address,
				Pincode:     order.User.Pincode,
			},
		}
		responses = append(responses, response)
	}

	return responses, nil
}
