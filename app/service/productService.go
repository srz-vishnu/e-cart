package service

import (
	"e-cart/app/dto"
	"e-cart/app/internal"
	"e-cart/pkg/e"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
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

	Product_id, err := s.productRepo.CreateAndUpsertProductDetail(args)
	if err != nil {
		return nil, e.NewError(e.ErrCreateProduct, "Failed to save product details", err)
	}
	log.Info().Msgf("Successfully added product details %d", Product_id)

	return &dto.CreateProductResponds{
		ProductID: Product_id,
	}, nil
}

func (s *ProductServiceImpl) ListAllProduct(r *http.Request) ([]*dto.CatagoryListResponse, error) {

	allCatagoryLists, err := s.productRepo.GetAllProducts()
	if err != nil {
		return nil, e.NewError(e.ErrValidateRequest, "error while listing all product items", err)
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
		//fmt.Printf("all product details: %v", catagorylists)
		fmt.Printf("Category ID: %d, Name: %s, Description: %s\n", prodlist.CatagoryID, prodlist.CatagoryName, prodlist.CatagoryName)

	}

	// Print each product's details directly
	for _, cat := range catagorylists {
		fmt.Printf("Category ID: %d, Name: %s, Description: %s\n", cat.CatagoryID, cat.CatagoryName, cat.Description)
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
		return nil, e.NewError(e.ErrCreateBook, "error while catagory detials", err)
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
		return nil, e.NewError(e.ErrCreateBook, "error while getting category details by name", err)
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
		return nil, e.NewError(e.ErrCreateBook, "error while getting all brands", err)
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

	// Print each product's details directly
	for _, brand := range brandLists {
		fmt.Printf(" Brand name: %s, Brand_Id: %d, price: %v, categoryID: %d  stockcount: %d\n", brand.BrandName, brand.BrandId, brand.Price, brand.CategoryID, brand.StockCount)
	}

	return brandLists, nil
}

// UpdateCategory updates the category name
func (s *ProductServiceImpl) UpdateCategory(r *http.Request) error {

	args := &dto.UpdateCategory{}

	err := args.Parse(r)
	if err != nil {
		return nil
	}

	//validation
	err = args.Validate()
	if err != nil {
		return e.NewError(e.ErrValidateRequest, "error while validating", err)
	}
	log.Info().Msg("Successfully completed parsing and validation of request body")

	if args.CategoryName == "" {
		return errors.New("category name cannot be empty")
	}
	return s.productRepo.UpdateCategory(args.CategoryID, args.CategoryName)
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
		return errors.New("brand name cannot be empty")
	}
	err = s.productRepo.UpdateBrand(args.BrandId, args.BrandName, args.Price)
	if err != nil {
		return e.NewError(e.ErrCreateBook, "no such id", err)
	}
	return nil
}
