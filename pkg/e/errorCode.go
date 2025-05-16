package e

// Application wide codes
const (
	ErrCodeAuto            = 0
	ErrCodeInternalService = 666
)

// 400 errors
const (
	// ErrInvalidRequest : when post body, query param, or path
	// param is invalid, or any post body validation error is encountered
	ErrInvalidRequest int = 400000 + iota

	// ErrDecodeRequestBody : error when decode the request body
	ErrDecodeRequestBody

	// ErrValidateRequest : error when validating the request
	ErrValidateRequest

	// ErrCreateProduct : error when creating product
	ErrCreateProduct

	// ErrCreateUser : error when creating user
	ErrCreateUser

	// ErrGetAuthorById : error when getting author by id
	ErrGetAuthorById

	// ErrUpdateAuthor : error when updating author
	ErrUpdateAuthor

	// ErrGetAllAuthorDetails : error to get all other details
	ErrGetAllAuthorDetails

	// ErrDeleteAuthor : error while deleting an author
	ErrDeleteAuthor

	// ErrBlockUser : error while blocking a user
	ErrBlockUser

	// ErrUnblockUser : error while unblocking a user
	ErrUnblockUser

	// ErrGetUserDetails : error while getting user details
	ErrGetUserDetails

	// ErrGetOrderHistory : error while getting order history
	ErrGetOrderHistory

	// ErrGetCartDetails : error while getting cart details
	ErrGetCartDetails

	// ErrUpdateCart : error while updating cart
	ErrUpdateCart

	// ErrPlaceOrder : error while placing order
	ErrPlaceOrder

	// ErrUpdateStock : error while updating stock
	ErrUpdateStock

	// ErrInsufficientStock : error when stock is insufficient
	ErrInsufficientStock

	// ErrGetFavoriteBrands : error while getting favorite brands
	ErrGetFavoriteBrands

	// ErrUpdateFavorites : error while updating favorites
	ErrUpdateFavorites

	// Product Service Errors
	// ErrListProducts : error while listing all products
	ErrListProducts

	// ErrGetCategory : error while getting category details
	ErrGetCategory

	// ErrGetBrand : error while getting fav brand details
	ErrGetFavBrand

	// ErrGetBrand : error while getting brand details
	ErrGetBrand

	// ErrUpdateCategory : error while updating category
	ErrUpdateCategory

	// ErrUpdateBrand : error while updating brand
	ErrUpdateBrand

	// User Service Errors
	// ErrLoginUser : error during user login
	ErrLoginUser

	// ErrInvalidCredentials : error when credentials are invalid
	ErrInvalidCredentials

	// ErrUserBlocked : error when user is blocked
	ErrUserBlocked

	// ErrGenerateToken : error while generating JWT token
	ErrGenerateToken

	// ErrAddToCart : error while adding item to cart
	ErrAddToCart

	// ErrClearCart : error while clearing cart
	ErrClearCart

	// ErrViewCart : error while viewing cart
	ErrViewCart

	// ErrAddToFavorites : error while adding to favorites
	ErrAddToFavorites

	// ErrGetFavorites : error while getting favorites
	ErrGetFavorites

	// ErrUpdateUserProfile : error while updating user profile
	ErrUpdateUserProfile
)

// 404 errors
const (
	// ErrResourceNotFound : when no record corresponding to the requested id is found in the DB
	ErrResourceNotFound int = 404000 + iota

	// ErrUserNotFound : when user is not found
	ErrUserNotFound

	// ErrProductNotFound : when product is not found
	ErrProductNotFound

	// ErrOrderNotFound : when order is not found
	ErrOrderNotFound

	// ErrCartNotFound : when cart is not found
	ErrCartNotFound

	// ErrCategoryNotFound : when category is not found
	ErrCategoryNotFound

	// ErrBrandNotFound : when brand is not found
	ErrBrandNotFound
)

// 500 errors
const (
	// ErrInternalServer : the default error, which is unexpected from the developers
	ErrInternalServer int = 500000 + iota

	// ErrExecuteSQL : when execute the sql, meet unexpected error
	ErrExecuteSQL

	// ErrDatabaseOperation : when database operation fails
	ErrDatabaseOperation

	// ErrContextError : when context related operations fail
	ErrContextError

	// ErrTransactionError : when database transaction fails
	ErrTransactionError
)
