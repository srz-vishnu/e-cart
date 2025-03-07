package service

import (
	"e-cart/app/dto"
	"e-cart/app/internal"
	"e-cart/pkg/e"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ProductService interface {
	CreateProduct(r *http.Request) (*dto.CreateProductResponds, error)
	ListAllProduct(r *http.Request) ([]*dto.CatagoryListResponse, error)
	GetCatagoryById(r *http.Request) (*dto.CreateCategoryDetailRequest, error)
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

func (s *ProductServiceImpl) GetCatagoryById(r *http.Request) (*dto.CreateCategoryDetailRequest, error) {
	args := &dto.SearchByCatagoryIdRequest{}

	//parsing
	err := args.Parse(r)
	if err != nil {
		log.Printf("Error parsing ID: %v\n", err)
		return nil, e.NewError(e.ErrDecodeRequestBody, "error while parsing", err)
	}

	fmt.Printf("print something  %d", args.CatagoryId)
	cat, err := s.productRepo.GetCategoryByID(args.CatagoryId)
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while catagory detials", err)
	}

	// Map the category and brand data to the DTO
	var response dto.CreateCategoryDetailRequest

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

func (s *ProductServiceImpl) ListAllBrands(r *http.Request) ([]*dto.BrandDetailResponse, error) {

	allBrandList, err := s.productRepo.GetAllBrands()
	if err != nil {
		return nil, e.NewError(e.ErrCreateBook, "error while getting all brands", err)
	}
	log.Info().Msgf("Successfully got all brand details %v \n", allBrandList)

	var brandLists []*dto.BrandDetailResponse

	for _, catBrand := range allBrandList {
		brandList := dto.BrandDetailResponse{
			BrandName:  catBrand.BrandName,
			Price:      catBrand.Price,
			StockCount: catBrand.StockCount,
		}
		brandLists = append(brandLists, &brandList)
		fmt.Printf("brand items %v", brandLists)
	}

	// Print each product's details directly
	for _, brand := range brandLists {
		fmt.Printf(" Brand name: %s, price: %v, stockcount: %d\n", brand.BrandName, brand.Price, brand.StockCount)
	}

	return brandLists, nil
}
