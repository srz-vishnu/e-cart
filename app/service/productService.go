package service

import (
	"e-cart/app/dto"
	"e-cart/app/internal"
	"e-cart/pkg/e"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ProductService interface {
	CreateProduct(r *http.Request) (*dto.CreateProductResponds, error)
}

type ProductServiceImpl struct {
	productRepo internal.ProductRepo
}

func NewProductService(productRepo internal.ProductRepo) ProductService {
	return &ProductServiceImpl{
		productRepo: productRepo,
	}
}

func (s *ProductServiceImpl) CreateProduct(r *http.Request) (*dto.CreateProductResponds, error) {
	args := &dto.CreateProductDetailRequest{}

	//parsing
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

	Product_id, err := s.productRepo.CreateProductDetail(args)
	if err != nil {
		return nil, e.NewError(e.ErrCreateProduct, "Failed to save product details", err)
	}
	log.Info().Msgf("Successfully added product details %d", Product_id)

	return &dto.CreateProductResponds{
		ProductID: Product_id,
	}, nil
}
