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

type ProductService interface {
	CreateProduct(r *http.Request) (*dto.CreateProductResponds, error)
	ListAllProduct(r *http.Request) ([]*dto.CatagoryListResponse, error)
	GetCatagoryById(r *http.Request) (*dto.CategoryDetailResponse, error)
	GetCatagoryByName(r *http.Request) (*dto.CategoryDetailResponse, error)
	ListAllBrands(r *http.Request) ([]*dto.BrandDetailResponse, error)
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
	args := &dto.CreateCategoryDetailRequest{}

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

	category, err := s.productRepo.CreateAndUpsertProductDetail(args)
	if err != nil {
		return nil, e.NewError(e.ErrCreateProduct, "Failed to save product details", err)
	}
	log.Info().Msgf("Successfully added product details %d", category.ID)

	var brands []dto.BrandResponse
	for _, b := range category.Brands {
		brands = append(brands, dto.BrandResponse{
			BrandName:  b.BrandName,
			Price:      b.Price,
			StockCount: b.StockCount,
			ImageLink:  b.ImageLink,
		})
	}

	return &dto.CreateProductResponds{
		ProductID:   category.ID,
		Category:    category.Categoryname,
		Description: category.Description,
		Brands:      brands,
	}, nil
}

func (s *ProductServiceImpl) ListAllProduct(r *http.Request) ([]*dto.CatagoryListResponse, error) {
	allCatagoryLists, err := s.productRepo.GetAllProducts()
	if err != nil {
		return nil, e.NewError(e.ErrListProducts, "error while listing all product items", err)
	}
	log.Info().Msgf("Successfully got all product details %v", allCatagoryLists)

	var catagorylists []*dto.CatagoryListResponse

	for _, pro := range allCatagoryLists {
		prodlist := dto.CatagoryListResponse{
			CatagoryID:   pro.ID,
			CatagoryName: pro.Categoryname,
			Description:  pro.Description,
		}
		catagorylists = append(catagorylists, &prodlist)
		fmt.Printf("Category ID: %d, Name: %s, Description: %s\n", prodlist.CatagoryID, prodlist.CatagoryName, prodlist.CatagoryName)
	}

	return catagorylists, nil
}

func (s *ProductServiceImpl) GetCatagoryById(r *http.Request) (*dto.CategoryDetailResponse, error) {
	args := &dto.SearchByCatagoryIdRequest{}

	//parsing
	err := args.Parse(r)
	if err != nil {
		log.Printf("Error parsing ID: %v\n", err)
		return nil, e.NewError(e.ErrDecodeRequestBody, "error while parsing", err)
	}

	cat, err := s.productRepo.GetCategoryByID(args.CatagoryId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewError(e.ErrCategoryNotFound, "category not found", err)
		}
		return nil, e.NewError(e.ErrGetCategory, "error while getting category details", err)
	}

	// Map the category and brand data to the DTO
	var response dto.CategoryDetailResponse

	response.CategoryID = cat.ID
	response.CategoryName = cat.Categoryname
	response.Description = cat.Description

	// Map the brands
	for _, brand := range cat.Brands {
		brandResp := dto.BrandDetailRequest{
			BrandName:  brand.BrandName,
			Price:      brand.Price,
			StockCount: brand.StockCount,
		}
		response.Brands = append(response.Brands, brandResp)
	}

	return &response, nil
}

func (s *ProductServiceImpl) GetCatagoryByName(r *http.Request) (*dto.CategoryDetailResponse, error) {
	args := &dto.SearchProductByNameRequest{}

	err := args.Parse(r)
	if err != nil {
		log.Printf("Error parsing category name: %v\n", err)
		return nil, e.NewError(e.ErrDecodeRequestBody, "error while parsing", err)
	}

	categoryDetails, err := s.productRepo.GetCategoryByName(args.CategoryName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewError(e.ErrCategoryNotFound, "category not found", err)
		}
		return nil, e.NewError(e.ErrGetCategory, "error while getting category details by name", err)
	}

	var response dto.CategoryDetailResponse

	response.CategoryID = categoryDetails.ID
	response.CategoryName = categoryDetails.Categoryname
	response.Description = categoryDetails.Description

	for _, brand := range categoryDetails.Brands {
		brandResp := dto.BrandDetailRequest{
			BrandName:  brand.BrandName,
			Price:      brand.Price,
			StockCount: brand.StockCount,
		}
		response.Brands = append(response.Brands, brandResp)
	}

	return &response, nil
}

func (s *ProductServiceImpl) ListAllBrands(r *http.Request) ([]*dto.BrandDetailResponse, error) {
	allBrandList, err := s.productRepo.GetAllBrands()
	if err != nil {
		return nil, e.NewError(e.ErrGetBrand, "error while getting all brands", err)
	}
	log.Info().Msgf("Successfully got all brand details %v \n", allBrandList)

	var brandLists []*dto.BrandDetailResponse

	for _, catBrand := range allBrandList {
		brandList := dto.BrandDetailResponse{
			BrandName:    catBrand.BrandName,
			BrandId:      catBrand.ID,
			Price:        catBrand.Price,
			StockCount:   catBrand.StockCount,
			CategoryID:   catBrand.CategoryID,
			CategoryName: catBrand.Category.Categoryname,
		}
		brandLists = append(brandLists, &brandList)
		fmt.Printf("brand items %v", brandLists)
	}

	return brandLists, nil
}

// UpdateCategory updates the category name
func (s *ProductServiceImpl) UpdateCategory(r *http.Request) error {
	args := &dto.UpdateCategory{}

	err := args.Parse(r)
	if err != nil {
		return e.NewError(e.ErrDecodeRequestBody, "error while parsing", err)
	}

	//validation
	err = args.Validate()
	if err != nil {
		return e.NewError(e.ErrValidateRequest, "error while validating", err)
	}
	log.Info().Msg("Successfully completed parsing and validation of request body")

	if args.CategoryName == "" {
		return e.NewError(e.ErrValidateRequest, "category name cannot be empty", nil)
	}

	err = s.productRepo.UpdateCategory(args.CategoryID, args.CategoryName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewError(e.ErrCategoryNotFound, "category not found", err)
		}
		return e.NewError(e.ErrUpdateCategory, "failed to update category", err)
	}
	return nil
}

// UpdateBrand updates the brand name
func (s *ProductServiceImpl) UpdateBrand(r *http.Request) error {
	args := &dto.UpdateBrand{}

	err := args.Parse(r)
	if err != nil {
		return e.NewError(e.ErrDecodeRequestBody, "error while parsing", err)
	}

	//validation
	err = args.Validate()
	if err != nil {
		return e.NewError(e.ErrValidateRequest, "error while validating", err)
	}
	log.Info().Msg("Successfully completed parsing and validation of request body")

	if args.BrandName == "" {
		return e.NewError(e.ErrValidateRequest, "brand name cannot be empty", nil)
	}

	err = s.productRepo.UpdateBrand(args.BrandId, args.BrandName, args.Price)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewError(e.ErrBrandNotFound, "brand not found", err)
		}
		return e.NewError(e.ErrUpdateBrand, "failed to update brand", err)
	}
	return nil
}
