package service

import (
	"e-cart/app/dto"
	"e-cart/app/internal"
	"e-cart/pkg/e"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserService interface {
	SaveUserDetails(r *http.Request) (*dto.SaveUserResponse, error)
	UpdateUserDetails(r *http.Request) error
	AddItemToCart(r *http.Request) (*dto.CartItemResponse, error)
	PlaceOrder(r *http.Request) (*dto.ItemOrderedResponse, error)
}

type userServiceImpl struct {
	userRepo internal.UserRepo
}

func NewUserService(userRepo internal.UserRepo) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

func (s *userServiceImpl) SaveUserDetails(r *http.Request) (*dto.SaveUserResponse, error) {
	args := &dto.UserDetailSaveRequest{}

	// parsing the req.body
	err := args.Parse(r)
	if err != nil {
		return nil, e.NewError(e.ErrDecodeRequestBody, "error while parsing", err)
	}

	//validation
	err = args.Validate()
	if err != nil {
		return nil, e.NewError(e.ErrValidateRequest, "error while validating", err)
	}
	log.Info().Msg("Successfully completed parsing and validation of request body")

	userID, err := s.userRepo.SaveUserDetails(args)
	if err != nil {
		return nil, e.NewError(e.ErrExecuteSQL, "error while creating user", err)
	}
	log.Info().Msgf("Successfully created user with id %d", userID)

	return &dto.SaveUserResponse{
		UserId: userID,
	}, nil
}

func (s *userServiceImpl) UpdateUserDetails(r *http.Request) error {

	args := &dto.UpdateUserDetailRequest{}

	err := args.Parse(r)
	if err != nil {
		return e.NewError(e.ErrDecodeRequestBody, "error while parsing", err)
	}

	err = args.Validate()
	if err != nil {
		return e.NewError(e.ErrValidateRequest, "error validating the req body", err)
	}
	log.Info().Msg("Successfully completed parsing and validation of request body")

	err = s.userRepo.UpdateUserDetails(args, args.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewError(e.ErrResourceNotFound, "user not found in the table", err)
		}
		return e.NewError(e.ErrCreateUser, "failed to update user details", err)
	}
	log.Info().Msg("Succesfully updated user details")

	return nil
}

func (s *userServiceImpl) AddItemToCart(r *http.Request) (*dto.CartItemResponse, error) {
	args := &dto.AddItemToCart{}

	err := args.Parse(r)
	if err != nil {
		return nil, e.NewError(e.ErrDecodeRequestBody, "errro while parsing", err)
	}

	err = args.Validate()
	if err != nil {
		return nil, e.NewError(e.ErrDecodeRequestBody, "error validating the req.body", err)
	}
	log.Info().Msg("Successfully completed parsing and validation of request body")

	// checking is user active or not
	isActive, err := s.userRepo.IsUserActive(args.UserID)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while checking user details", err)
	}

	if !isActive {
		log.Info().Msg("User is not active.")
		return nil, e.NewError(e.ErrCreateBook, "user is blocked or inactive", nil)
	}
	log.Info().Msg("User is active")

	// getting product price
	prodDetails, err := s.userRepo.GetProductDetails(args.ProductID, args.Quantity)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while getting product details", err)
	}
	log.Info().Msgf("succesfully got product details of product id %v", prodDetails.Price)

	// checking product already exsist in cart, if not adding those items
	cartData, err := s.userRepo.AddOrUpdateCart(args.UserID, prodDetails, args.Quantity)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while adding items to the cart", err)
	}
	// Calculate total price
	totalPrice := float64(cartData.Quantity) * cartData.Price
	//brand := prodDetails.BrandName

	// Create the cart item response
	cartItemResponse := dto.CartItemResponse{
		UserID:     cartData.UserID,
		ProductID:  cartData.ProductID,
		Quantity:   cartData.Quantity,
		Price:      cartData.Price,
		BrandName:  cartData.Product.BrandName,
		TotalPrice: totalPrice,
	}
	log.Info().Msgf("Cart Item Response: %+v", cartItemResponse)

	return &cartItemResponse, nil
}

func (s userServiceImpl) PlaceOrder(r *http.Request) (*dto.ItemOrderedResponse, error) {

	args := dto.PlaceOrderFromCart{}

	err := args.Parse(r)
	if err != nil {
		return nil, e.NewError(e.ErrDecodeRequestBody, "errro while parsing", err)
	}

	err = args.Validate()
	if err != nil {
		return nil, e.NewError(e.ErrDecodeRequestBody, "error validating the req.body", err)
	}
	log.Info().Msg("Successfully completed parsing and validation of request body")

	// checking is user active or not
	isActive, err := s.userRepo.IsUserActive(args.UserID)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while checking user details", err)
	}

	if !isActive {
		log.Info().Msg("User is not active.")
		return nil, e.NewError(e.ErrCreateBook, "user is blocked or inactive", nil)
	}
	log.Info().Msg("User is active")

	// fetch cartItems against userId
	cartItems, err := s.userRepo.FetchCartItems(args.UserID)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while fetching cart details", err)
	}

	// Log details of each cart item
	for _, cartItem := range cartItems {
		log.Info().Msgf("CartID: %d, ProductID: %d, Quantity: %d, Price: %.2f, BrandName: %s", cartItem.ID, cartItem.ProductID, cartItem.Quantity,
			cartItem.Price, cartItem.Product.BrandName)
	}

	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += item.Price * float64(item.Quantity)
	}
	log.Info().Msgf("totalAmount is %v :", totalAmount)

	newOrder, err := s.userRepo.CreateOrder(args.UserID, totalAmount, cartItems)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while creating order", err)
	}
	log.Info().Msgf("Order ID: %d, Total: %.2f, UserID: %d", newOrder.ID, newOrder.Total, newOrder.UserID)

	items, err := s.userRepo.CreateOrderItems(newOrder.ID, cartItems)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while create order items", err)
	}
	// Log each item
	for _, item := range items {
		log.Info().Msgf("Item ID: %d, OrderID: %d, ProductID: %d, BrandName: %s, Quantity: %d, Price: %.2f",
			item.ID, item.OrderID, item.ProductID, item.Product.BrandName, item.Quantity, item.Price)
	}

	log.Info().Msg("Success")

	itemOrderedResponse := dto.ItemOrderedResponse{
		OrderID:    items[0].OrderID,
		ProductID:  items[0].ProductID,
		Quantity:   items[0].Quantity,
		CatagoryID: items[0].Product.CategoryID,
		BrandName:  items[0].Product.BrandName,
		TotalPrice: totalAmount,
	}

	return &itemOrderedResponse, nil
}
