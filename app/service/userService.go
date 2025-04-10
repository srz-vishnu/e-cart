package service

import (
	"e-cart/app/dto"
	helper "e-cart/app/helper"
	"e-cart/app/internal"
	"e-cart/pkg/e"
	"e-cart/pkg/jwt"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserService interface {
	SaveUserDetails(r *http.Request) (*dto.SaveUserResponse, error)
	LoginUser(r *http.Request) (*dto.LoginResponse, error)
	UpdateUserDetails(r *http.Request) error
	ViewUserCart(r *http.Request) ([]*dto.ViewCart, error)
	ClearCart(r *http.Request) error
	AddItemToCart(r *http.Request) (*dto.CartItemResponse, error)
	OrderHistory(r *http.Request) (*dto.ItemOrderedResponse, error)
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

func (s *userServiceImpl) LoginUser(r *http.Request) (*dto.LoginResponse, error) {

	args := &dto.LoginRequest{}

	// parsing the req.body
	err := args.Parse(r)
	if err != nil {
		return nil, e.NewError(e.ErrValidateRequest, "error while parsing", err)
	}

	//validation
	err = args.Validate()
	if err != nil {
		return nil, e.NewError(e.ErrValidateRequest, "error while vlidating", err)
	}
	log.Info().Msg("Successfully completed parsing and validation of request body")

	// Fetching user from database
	user, err := s.userRepo.GetUserByUsername(args.Username)
	if err != nil {
		return nil, e.NewError(e.ErrResourceNotFound, "user not found", err)
	}

	// Check if user is nil
	if user == nil {
		return nil, e.NewError(e.ErrResourceNotFound, "user not found", err)
	}
	log.Info().Msgf("the user is %s", user.Username)

	// Validate password
	if user.Password != args.Password {
		err := fmt.Errorf("invalid password for user %s", user.Username)
		return nil, e.NewError(e.ErrDecodeRequestBody, "invalid password", err)
	}

	// Generating JWT Token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.IsAdmin) //userid and username in the token
	if err != nil {
		return nil, e.NewError(e.ErrInternalServer, "failed to generate token", err)
	}
	fmt.Printf("the token is %s : \n ", token)

	return &dto.LoginResponse{
		Token: token,
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
	prodDetails, err := s.userRepo.GetProductDetails(args.CategoryID, args.BrandId)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while getting product details", err)
	}
	log.Info().Msgf("succesfully got product details of barnd name %s price %f stock %d", prodDetails.BrandName, prodDetails.Price, prodDetails.StockCount)

	// Check if the requested quantity exceeds available stock
	if args.Quantity > prodDetails.StockCount {
		log.Info().Msgf("Only %d units of %s are available, but you requested %d", prodDetails.StockCount, prodDetails.BrandName, args.Quantity)
		return nil, e.NewError(e.ErrInternalServer, "dont have enough stocks to meet your requested stock count", errors.New("stock insufficient"))
	}
	log.Info().Msg("Requested quantity is available")

	totalAmount := prodDetails.Price * float64(args.Quantity)
	log.Info().Msgf("totalAmount is %v :", totalAmount)

	// checking product already exsist in cart, if not adding those items
	err = s.userRepo.AddOrUpdateCart(args.UserID, prodDetails, args.Quantity, totalAmount)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while adding items to the cart", err)
	}
	log.Info().Msg("Succesfully added items to the cart")

	// Get the updated cart details along with product info
	cartData, err := s.userRepo.GetCartWithProductDetails(args.UserID, prodDetails.ID)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while retrieving cart with product details", err)
	}
	log.Info().Msgf("succesfully got cart details of barnd name %s ,price %f ,stock %d", cartData.Brand.BrandName, cartData.Brand.Price, cartData.Brand.StockCount)

	// Create the cart item response
	cartItemResponse := dto.CartItemResponse{
		UserID:     cartData.UserID,
		ProductID:  cartData.ProductID,
		Quantity:   cartData.Quantity,
		Price:      cartData.Price,
		BrandName:  cartData.Brand.BrandName,
		TotalPrice: totalAmount,
	}
	log.Info().Msgf("Cart Item Response: %+v", cartItemResponse)

	return &cartItemResponse, nil

}

func (s *userServiceImpl) ViewUserCart(r *http.Request) ([]*dto.ViewCart, error) {

	UserId, err := helper.GetUserIDFromContext(r.Context())
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while getting userId from ctx", err)
	}
	log.Info().Msgf("userId of the user logged in %d", UserId)

	isActive, err := s.userRepo.IsUserActive(UserId)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while checking user details", err)
	}

	if !isActive {
		log.Info().Msg("User is not active.")
		return nil, e.NewError(e.ErrCreateBook, "user is blocked or inactive", nil)
	}
	log.Info().Msg("User is active")

	cartDetails, err := s.userRepo.ViewCart(UserId)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "not able to see the cart associated with the user", err)
	}

	// only for formatting - for log view

	log.Info().Msg("User cart list:")
	for _, cart := range cartDetails {
		log.Info().Msgf("ProductID: %d, Quantity: %d, Price: %.2f, Brand: %s",
			cart.ProductID, cart.Quantity, cart.Price, cart.Brand.BrandName)
	}
	/////////////////////////////////////////////////////////

	var cartlists []*dto.ViewCart

	for _, carts := range cartDetails {
		list := dto.ViewCart{
			ProductID:   carts.ProductID,
			BrandName:   carts.Brand.BrandName,
			Quantity:    carts.Quantity,
			Price:       carts.Price,
			TotalAmount: carts.Price * float64(carts.Quantity),
		}

		cartlists = append(cartlists, &list)
	}

	return cartlists, nil
}

func (s *userServiceImpl) ClearCart(r *http.Request) error {

	UserId, err := helper.GetUserIDFromContext(r.Context())
	if err != nil {
		return e.NewError(e.ErrCreateBook, "error while getting userId from ctx", err)
	}
	log.Info().Msgf("userId of the user logged in %d", UserId)

	isActive, err := s.userRepo.IsUserActive(UserId)
	if err != nil {
		return e.NewError(e.ErrCreateBook, "error while checking user details", err)
	}

	if !isActive {
		log.Info().Msg("User is not active.")
		return e.NewError(e.ErrCreateBook, "user is blocked or inactive", nil)
	}
	log.Info().Msg("User is active")

	err = s.userRepo.ClearCart(UserId)
	if err != nil {
		return e.NewError(e.ErrCreateBook, "failed to clear cart", err)
	}
	log.Info().Msg("Cart cleared succesfully")

	return nil
}

func (s *userServiceImpl) PlaceOrder(r *http.Request) (*dto.ItemOrderedResponse, error) {
	args := dto.PlaceOrderFromCart{}

	// Parse and validate request
	err := args.Parse(r)
	if err != nil {
		return nil, e.NewError(e.ErrDecodeRequestBody, "error while parsing", err)
	}

	err = args.Validate()
	if err != nil {
		return nil, e.NewError(e.ErrDecodeRequestBody, "error validating the req.body", err)
	}
	log.Info().Msg("Successfully completed parsing and validation of request body")

	// Check if user is active
	isActive, err := s.userRepo.IsUserActive(args.UserID)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while checking user details", err)
	}

	if !isActive {
		log.Info().Msg("User is not active.")
		return nil, e.NewError(e.ErrCreateBook, "user is blocked or inactive", nil)
	}
	log.Info().Msg("User is active")

	// Fetch cart items
	cartItems, err := s.userRepo.FetchCartItems(args.UserID, args.CartID)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while fetching cart details", err)
	}

	// Log cart items
	for _, cartItem := range cartItems {
		log.Info().Msgf("CartID: %d, ProductID: %d, Quantity: %d, Price: %.2f, BrandName: %s",
			cartItem.ID, cartItem.ProductID, cartItem.Quantity, cartItem.Price, cartItem.Brand.BrandName)
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += item.Price * float64(item.Quantity)
	}
	log.Info().Msgf("totalAmount is %v :", totalAmount)

	// Create order and items in a single transaction
	newOrder, orderItems, err := s.userRepo.CreateOrder(args.UserID, totalAmount, cartItems)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while creating order", err)
	}
	log.Info().Msgf("Order ID: %d, Total: %.2f, UserID: %d", newOrder.ID, newOrder.Total, newOrder.UserID)

	// Get user details for response
	user, err := s.userRepo.GetUserByID(args.UserID)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while fetching user details", err)
	}

	// Update stock count
	updatedBrands, err := s.userRepo.UpdateStockCount(orderItems)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while updating stock count", err)
	}

	// Log updated stock details
	for _, brand := range updatedBrands {
		log.Info().Msgf("Product ID: %d, Brand Name: %s, New Stock Count: %d", brand.ID, brand.BrandName, brand.StockCount)
	}
	log.Info().Msg("Successfully updated stock count after order placement")

	// Update cart status
	err = s.userRepo.UpdateCartOrderStatus(args.UserID, newOrder.ID, args.CartID)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while updating cart status to not active", err)
	}
	log.Info().Msg("Successfully updated the cart status to false after order being placed")

	// Build response
	itemOrderedResponse := dto.ItemOrderedResponse{
		OrderID:    newOrder.ID,
		TotalPrice: totalAmount,
		UserDetails: dto.UserDetailsResponse{
			Username:    user.Username,
			Address:     user.Address,
			Pincode:     user.Pincode,
			PhoneNumber: user.Phonenumber,
			Email:       user.Mail,
		},
		Items: make([]dto.OrderItemResponse, 0, len(orderItems)),
	}

	for _, item := range orderItems {
		itemOrderedResponse.Items = append(itemOrderedResponse.Items, dto.OrderItemResponse{
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			CategoryID: item.Product.CategoryID,
			BrandName:  item.Product.BrandName,
			Price:      item.Price,
		})
	}

	return &itemOrderedResponse, nil
}

func (s *userServiceImpl) OrderHistory(r *http.Request) (*dto.ItemOrderedResponse, error) {
	return nil, nil
}
