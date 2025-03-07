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
}

type AdminServiceImpl struct {
	adminRepo internal.AdminRepo
}

func NewAdminService(adminRepo internal.AdminRepo) AdminService {
	return &AdminServiceImpl{
		adminRepo: adminRepo,
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
			return e.NewError(e.ErrResourceNotFound, "user not found in the table", err)
		}
		return e.NewError(e.ErrCreateUser, "failed to block user", err)
	}
	log.Info().Msg("user with ID %d has been successfully blockd")

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
			return e.NewError(e.ErrResourceNotFound, "user not found in the table", err)
		}
		return e.NewError(e.ErrCreateUser, "failed to unblock user", err)
	}
	log.Info().Msg("user with ID %d has been successfully unblockd")

	return nil
}

func (s *AdminServiceImpl) GetAllUserDetail(r *http.Request) ([]*dto.AllUserDetails, error) {

	allUserDetails, err := s.adminRepo.GetAllUsers()
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "failed to get all user details", err)
	}

	var userDetails []*dto.AllUserDetails

	for _, user := range allUserDetails {
		userDetail := &dto.AllUserDetails{
			UserName: user.Username,
			Mail:     user.Mail,
			Address:  user.Address,
			Phone:    user.Phonenumber,
		}

		userDetails := append(userDetails, userDetail)
		fmt.Printf("all user details: %v", userDetails)
	}

	return userDetails, nil
}
