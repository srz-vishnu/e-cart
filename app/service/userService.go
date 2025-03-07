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
	AddItemToCart(r *http.Request) error
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

func (s *userServiceImpl) AddItemToCart(r *http.Request) error {
	args := &dto.AddItemToCart{}

	err := args.Parse(r)
	if err != nil {
		return e.NewError(e.ErrDecodeRequestBody, "errro while parsing", err)
	}

	err = args.Validate()
	if err != nil {
		return e.NewError(e.ErrDecodeRequestBody, "error validating the req.body", err)
	}
	log.Info().Msg("Successfully completed parsing and validation of request body")

	err = s.userRepo.AddItemToCart(args)
	if err != nil {
		return e.NewError(e.ErrCreateBook, "failed to add items to the cart", err)
	}
	log.Info().Msg("Succesfully added  items to cart")

	return nil
}
