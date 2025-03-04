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

	//ErrCreateBook : error while creating book
	ErrCreateBook
)

// 404 errors
const (
	// ErrResourceNotFound : when no record corresponding to the requested id is found in the DB
	ErrResourceNotFound int = 404000 + iota
)

// 500 errors
const (
	// ErrInternalServer : the default error, which is unexpected from the developers
	ErrInternalServer int = 500000 + iota

	// ErrExecuteSQL : when execute the sql, meet unexpected error
	ErrExecuteSQL
)
