package service

import (
	"context"
	"e-cart/app/dto"
	helper "e-cart/app/helper"
	"e-cart/app/internal"
	"e-cart/pkg/e"
	"e-cart/pkg/jwt"
	hash "e-cart/pkg/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserService interface {
	SaveUserDetails(r *http.Request) (*dto.SaveUserResponse, error)
	LoginUser(r *http.Request) (*dto.LoginResponse, error)
	UpdateUserDetails(r *http.Request) error
	GetUserDetails(r *http.Request) (*dto.GetUserDetailsResponse, error)
	ChangePassword(r *http.Request) error
	ViewUserCart(r *http.Request) ([]*dto.ViewCart, error)
	ClearCart(r *http.Request) error
	AddItemToCart(r *http.Request) (*dto.CartItemResponse, error)
	PlaceOrder(r *http.Request) (*dto.ItemOrderedResponse, error)
	OrderHistory(r *http.Request) ([]*dto.ItemOrderedResponse, error)
	AddItemsToFavourites(r *http.Request) error
	GetUserFavouriteBrands(r *http.Request) ([]dto.FavoriteBrandResponse, error)
}

type userServiceImpl struct {
	userRepo      internal.UserRepo
	contextHelper helper.ContextHelper
	bcryptPackage hash.BcryptPackage
}

func NewUserService(userRepo internal.UserRepo, ctxHelper helper.ContextHelper, hashPassword hash.BcryptPackage) UserService {
	return &userServiceImpl{
		userRepo:      userRepo,
		contextHelper: ctxHelper,
		bcryptPackage: hashPassword,
	}
}

func (s *userServiceImpl) getUserIDAndCheckStatus(ctx context.Context) (int64, error) {
	userID, err := s.contextHelper.GetUserID(ctx)
	if err != nil {
		return 0, e.NewError(e.ErrContextError, "error while getting userId from ctx", err)
	}
	log.Info().Msgf("userId of the user logged in %d", userID)

	isActive, err := s.userRepo.IsUserActive(userID)
	if err != nil {
		return 0, e.NewError(e.ErrGetUserDetails, "error while checking user details", err)
	}

	if !isActive {
		log.Info().Msg("User is not active.")
		return 0, e.NewError(e.ErrUserBlocked, "user is blocked or inactive", nil)
	}
	log.Info().Msg("User is active")

	return userID, nil
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

	// Check if username already exists
	existingUser, err := s.userRepo.GetUserByUsername(args.UserName)
	if err == nil && existingUser != nil {
		log.Info().Msgf("Username %s is already exist", args.UserName)
		return nil, e.NewError(e.ErrUserNameAlreadyExists, "username already exists", nil)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Only return error if it's not 'record not found'
		return nil, e.NewError(e.ErrInternal, "error checking existing user", err)
	}

	// hashing the password for protection
	hashedPassword, err := s.bcryptPackage.HashPassword(args.Password)
	if err != nil {
		return nil, e.NewError(e.ErrHashPassword, "failed to hash password", err)
	}
	args.Password = hashedPassword

	userID, err := s.userRepo.SaveUserDetails(args)
	if err != nil {
		return nil, e.NewError(e.ErrCreateUser, "error while creating user", err)
	}
	log.Info().Msgf("Successfully created user with id %d", userID)

	return &dto.SaveUserResponse{
		UserId: userID,
	}, nil
}

// func (s *userServiceImpl) GetUserDetails(r *http.Request) (*dto.GetUserDetailsResponse, error) {
// 	strID := chi.URLParam(r, "userid")
// 	intID, err := strconv.Atoi(strID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	UserID := int64(intID)

// 	userDetails, err := s.userRepo.GetUserDetailByID(UserID)
// 	if err != nil {
// 		return nil, e.NewError(e.ErrGetUserDetails, "failed to get user details", err)
// 	}

// 	if userDetails.IsAdmin {
// 		log.Info().Msg("the user is an admin")
// 	} else {
// 		log.Info().Msg("the user is a regular user")
// 	}

// 	// Check if user is active
// 	if !userDetails.Status {
// 		err := fmt.Errorf("user %s is blocked", userDetails.Username)
// 		return nil, e.NewError(e.ErrUserBlocked, "user is blocked", err)
// 	}

// 	response := &dto.GetUserDetailsResponse{
// 		UserName: userDetails.Username,
// 		Mail:     userDetails.Mail,
// 		Phone:    userDetails.Phonenumber,
// 		Address:  userDetails.Address,
// 		Pincode:  userDetails.Pincode,
// 	}

// 	return response, nil
// }

func (s *userServiceImpl) ChangePassword(r *http.Request) error {
	args := &dto.ChangePasswordRequest{}

	userID, err := s.getUserIDAndCheckStatus(r.Context())
	if err != nil {
		return nil
	}

	err = args.Parse(r)
	if err != nil {
		return e.NewError(e.ErrDecodeRequestBody, "error while parsing", err)
	}

	err = args.Validate()
	if err != nil {
		return e.NewError(e.ErrValidateRequest, "error validating the req.body", err)
	}
	log.Info().Msg("Successfully completed parsing and validation of request body")

	userDetails, err := s.userRepo.GetUserDetailByID(userID)
	if err != nil {
		return e.NewError(e.ErrGetUserDetails, "error while fetching user details", err)
	}

	// Validate password
	if args.NewPassword != args.ConfirmPassword {
		err := fmt.Errorf("mismatching password for user %s", userDetails.Username)
		return e.NewError(e.ErrMismatchingPassword, "mismatching new password and current password", err)
	}

	// comparing the current hashed_pwd from db with pwd from front_end
	passwordMatch := s.bcryptPackage.ComparePassword(userDetails.Password, args.CurrentPassword)
	if !passwordMatch {
		return e.NewError(e.ErrInvalidCredentials, "invalid password", nil)
	}

	hashPassword, err := s.bcryptPackage.HashPassword(args.NewPassword)
	if err != nil {
		return e.NewError(e.ErrHashPassword, "failed to hash password", err)
	}

	args.NewPassword = hashPassword

	err = s.userRepo.ChangePassword(userID, args.NewPassword)
	if err != nil {
		return e.NewError(e.ErrHashPassword, "failed to change the new password", err)
	}
	log.Info().Msg("Successfully updated password")

	return nil
}

func (s *userServiceImpl) GetUserDetails(r *http.Request) (*dto.GetUserDetailsResponse, error) {
	strID := chi.URLParam(r, "userid")
	intID, err := strconv.Atoi(strID)
	if err != nil {
		return nil, err
	}
	userID := int64(intID)

	// userID, err := s.getUserIDAndCheckStatus(r.Context())
	// if err != nil {
	// 	return nil, err
	// }

	userDetails, err := s.userRepo.GetUserDetailByID(userID)
	if err != nil {
		return nil, e.NewError(e.ErrGetUserDetails, "failed to get user details", err)
	}

	if userDetails.IsAdmin {
		log.Info().Msg("the user is an admin")
	} else {
		log.Info().Msg("the user is a regular user")
	}

	// Check if user is active
	if !userDetails.Status {
		err := fmt.Errorf("user %s is blocked", userDetails.Username)
		return nil, e.NewError(e.ErrUserBlocked, "user is blocked", err)
	}

	response := &dto.GetUserDetailsResponse{
		UserName: userDetails.Username,
		Mail:     userDetails.Mail,
		Address:  userDetails.Address,
		Pincode:  userDetails.Pincode,
		Phone:    userDetails.Phonenumber,
	}

	return response, nil
}

func (s *userServiceImpl) LoginUser(r *http.Request) (*dto.LoginResponse, error) {
	args := &dto.LoginRequest{}

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

	// Fetching user from database
	user, err := s.userRepo.GetUserByUsername(args.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewError(e.ErrUserNotFound, "user not found", err)
		}
		return nil, e.NewError(e.ErrLoginUser, "error during login", err)
	}

	// Check if user is nil
	if user == nil {
		return nil, e.NewError(e.ErrUserNotFound, "user not found", err)
	}

	if user.IsAdmin {
		log.Info().Msg("the user is an admin")
	} else {
		log.Info().Msg("the user is a regular user")
	}

	// Validating password using comparePassword
	// passwordMatch := s.bcryptPackage.ComparePassword(user.Password, args.Password)
	// if !passwordMatch {
	// 	return nil, e.NewError(e.ErrInvalidCredentials, "invalid password", nil)
	// }

	passwordMatch := s.bcryptPackage.ComparePassword(user.Password, args.Password)
	if !passwordMatch {
		log.Info().Msg("user is old, pwd is not hashed yet, so hash the old pwd and save it to db")
		if user.Password == args.Password {
			// Plaintext matched, so upgrade it to hashed password
			hashedPwd, err := s.bcryptPackage.HashPassword(args.Password)
			if err == nil {
				_ = s.userRepo.ChangePassword(user.ID, hashedPwd)
				log.Info().Msgf("Upgraded password to hashed format for user %s", user.Username)
			} else {
				//log.Error().Err(err).Msg("Failed to upgrade password hash")
				return nil, e.NewError(e.ErrHashPassword, "failed to hash password", err)
			}
		} else {
			// Neither match
			return nil, e.NewError(e.ErrInvalidCredentials, "invalid password", nil)
		}
	}

	// Check if user is active
	if !user.Status {
		err := fmt.Errorf("user %s is blocked", user.Username)
		return nil, e.NewError(e.ErrUserBlocked, "user is blocked", err)
	}

	// Generating JWT Token with isAdmin from database
	token, expiry, err := jwt.GenerateToken(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		return nil, e.NewError(e.ErrGenerateToken, "failed to generate token", err)
	}
	log.Info().Msgf("Generated token for user %s (Admin: %v)", user.Username, user.IsAdmin)

	// Saving generated token details on table
	err = s.userRepo.SaveToken(user.ID, token, expiry)
	if err != nil {
		return nil, e.NewError(e.ErrAddToFavorites, "failed to store login token", err)
	}
	log.Info().Msg("Generated token saved successfully")

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

	userID, err := s.getUserIDAndCheckStatus(r.Context())
	if err != nil {
		return err
	}

	// Check if username already exists
	existingUser, err := s.userRepo.GetUserByUsername(args.UserName)
	if err == nil && existingUser != nil && existingUser.ID != userID {
		return e.NewError(e.ErrUserNameAlreadyExists, "username already exists", nil)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Only return error if it's not 'record not found'
		return e.NewError(e.ErrInternal, "error checking existing user", err)
	}

	err = s.userRepo.UpdateUserDetails(args, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewError(e.ErrUserNotFound, "user not found in the table", err)
		}
		return e.NewError(e.ErrUpdateUserProfile, "failed to update user details", err)
	}
	log.Info().Msg("Successfully updated user details")

	return nil
}

func (s *userServiceImpl) AddItemToCart(r *http.Request) (*dto.CartItemResponse, error) {
	args := &dto.AddItemToCart{}

	userID, err := s.getUserIDAndCheckStatus(r.Context())
	if err != nil {
		return nil, err
	}

	err = args.Parse(r)
	if err != nil {
		return nil, e.NewError(e.ErrDecodeRequestBody, "error while parsing", err)
	}

	err = args.Validate()
	if err != nil {
		return nil, e.NewError(e.ErrValidateRequest, "error validating the req.body", err)
	}
	log.Info().Msg("Successfully completed parsing and validation of request body")

	// getting product price
	prodDetails, err := s.userRepo.GetProductDetails(args.CategoryID, args.BrandId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewError(e.ErrProductNotFound, "product not found", err)
		}
		return nil, e.NewError(e.ErrGetBrand, "error while getting product details", err)
	}
	log.Info().Msgf("successfully got product details of brand name %s price %f stock %d", prodDetails.BrandName, prodDetails.Price, prodDetails.StockCount)

	// Check if the requested quantity exceeds available stock
	if args.Quantity > prodDetails.StockCount {
		log.Info().Msgf("Only %d units of %s are available, but you requested %d", prodDetails.StockCount, prodDetails.BrandName, args.Quantity)
		return nil, e.NewError(e.ErrInsufficientStock, "insufficient stock available", errors.New("stock insufficient"))
	}
	log.Info().Msg("Requested quantity is available")

	totalAmount := prodDetails.Price * float64(args.Quantity)
	log.Info().Msgf("totalAmount is %v :", totalAmount)

	// checking product already exist in cart, if not adding those items
	err = s.userRepo.AddOrUpdateCart(userID, prodDetails, args.Quantity, totalAmount)
	if err != nil {
		return nil, e.NewError(e.ErrAddToCart, "error while adding items to the cart", err)
	}
	log.Info().Msg("Successfully added items to the cart")

	// Get the updated cart details along with product info
	cartData, err := s.userRepo.GetCartWithProductDetails(userID, prodDetails.ID)
	if err != nil {
		return nil, e.NewError(e.ErrGetCartDetails, "error while retrieving cart with product details", err)
	}
	log.Info().Msgf("successfully got cart details of brand name %s, price %f, stock %d", cartData.Brand.BrandName, cartData.Brand.Price, cartData.Brand.StockCount)

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
	userID, err := s.getUserIDAndCheckStatus(r.Context())
	if err != nil {
		return nil, err
	}

	cartDetails, err := s.userRepo.ViewCart(userID)
	if err != nil {
		return nil, e.NewError(e.ErrViewCart, "not able to see the cart associated with the user", err)
	}

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
	userID, err := s.getUserIDAndCheckStatus(r.Context())
	if err != nil {
		return err
	}

	err = s.userRepo.ClearCart(userID)
	if err != nil {
		return e.NewError(e.ErrClearCart, "failed to clear cart", err)
	}
	log.Info().Msg("Cart cleared successfully")

	return nil
}

func (s *userServiceImpl) PlaceOrder(r *http.Request) (*dto.ItemOrderedResponse, error) {
	userID, err := s.getUserIDAndCheckStatus(r.Context())
	if err != nil {
		return nil, err
	}

	args := dto.PlaceOrderFromCart{}

	// Parse and validate request
	err = args.Parse(r)
	if err != nil {
		return nil, e.NewError(e.ErrDecodeRequestBody, "error while parsing", err)
	}

	err = args.Validate()
	if err != nil {
		return nil, e.NewError(e.ErrValidateRequest, "error validating the req.body", err)
	}
	log.Info().Msg("Successfully completed parsing and validation of request body")

	// Fetch cart items
	cartItems, err := s.userRepo.FetchCartItems(userID, args.CartID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewError(e.ErrCartNotFound, "cart not found", err)
		}
		return nil, e.NewError(e.ErrGetCartDetails, "error while fetching cart details", err)
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += item.Price * float64(item.Quantity)
	}
	log.Info().Msgf("totalAmount is %v :", totalAmount)

	// Create order and items in a single transaction
	newOrder, orderItems, err := s.userRepo.CreateOrder(userID, totalAmount, cartItems)
	if err != nil {
		return nil, e.NewError(e.ErrPlaceOrder, "error while creating order", err)
	}
	log.Info().Msgf("Order ID: %d, Total: %.2f, UserID: %d", newOrder.ID, newOrder.Total, newOrder.UserID)

	// Get user details for response
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, e.NewError(e.ErrGetUserDetails, "error while fetching user details", err)
	}

	// Update stock count
	_, err = s.userRepo.UpdateStockCount(orderItems)
	if err != nil {
		return nil, e.NewError(e.ErrUpdateStock, "error while updating stock count", err)
	}

	// Update cart status
	err = s.userRepo.UpdateCartOrderStatus(userID, newOrder.ID, args.CartID)
	if err != nil {
		return nil, e.NewError(e.ErrUpdateCart, "error while updating cart status", err)
	}
	log.Info().Msg("Successfully updated the cart status to false after order being placed")

	// Build response
	itemOrderedResponse := dto.ItemOrderedResponse{
		OrderID: newOrder.ID,
		UserDetails: dto.UserDetailsResponse{
			Username:    user.Username,
			Address:     user.Address,
			Pincode:     user.Pincode,
			PhoneNumber: user.Phonenumber,
			Email:       user.Mail,
		},
		TotalPrice: totalAmount,
		Items:      make([]dto.OrderItemResponse, 0, len(orderItems)),
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

func (s *userServiceImpl) OrderHistory(r *http.Request) ([]*dto.ItemOrderedResponse, error) {
	userID, err := s.getUserIDAndCheckStatus(r.Context())
	if err != nil {
		return nil, err
	}

	userDetails, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, e.NewError(e.ErrGetUserDetails, "failed to get user details", err)
	}

	orderHistory, err := s.userRepo.GetOrderHistoryByUserID(userID)
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

func (s *userServiceImpl) AddItemsToFavourites(r *http.Request) error {
	args := dto.UserFavoriteBrandRequest{}

	// Parse and validate request
	err := args.Parse(r)
	if err != nil {
		return e.NewError(e.ErrDecodeRequestBody, "error while parsing", err)
	}

	err = args.Validate()
	if err != nil {
		return e.NewError(e.ErrValidateRequest, "error validating the req.body", err)
	}
	log.Info().Msg("Successfully completed parsing and validation of request body")

	userID, err := s.getUserIDAndCheckStatus(r.Context())
	if err != nil {
		return err
	}

	err = s.userRepo.AddOrUpdateFavorite(userID, args)
	if err != nil {
		return e.NewError(e.ErrAddToFavorites, "failed to update brand to the favourite list", err)
	}
	log.Info().Msg("Successfully updated the brand to the user favourite list")

	return nil
}

func (s *userServiceImpl) GetUserFavouriteBrands(r *http.Request) ([]dto.FavoriteBrandResponse, error) {
	userID, err := s.getUserIDAndCheckStatus(r.Context())
	if err != nil {
		return nil, err
	}

	brandIDs, err := s.userRepo.GetFavoriteBrandIDs(userID)
	if err != nil {
		return nil, e.NewError(e.ErrGetFavorites, "failed to get favorite brand IDs", err)
	}
	log.Info().Interface("favorite_brand_ids", brandIDs).Msg("Fetched favorite brand IDs")

	brands, err := s.userRepo.GetBrandsByIDs(brandIDs)
	if err != nil {
		return nil, e.NewError(e.ErrGetFavBrand, "failed to get favorite brand", err)
	}
	log.Info().Msg("Successfully got favorite brands")

	var resp []dto.FavoriteBrandResponse
	for _, b := range brands {
		resp = append(resp, dto.FavoriteBrandResponse{
			BrandID:   b.ID,
			BrandName: b.BrandName,
			Price:     b.Price,
			Stock:     b.StockCount,
			ImageLink: b.ImageLink,
		})
	}

	return resp, nil
}
