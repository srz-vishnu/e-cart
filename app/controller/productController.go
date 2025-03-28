package controller

import (
	"e-cart/app/service"
	"e-cart/pkg/api"
	"e-cart/pkg/e"
	"net/http"
)

type ProductController interface {
	CreateProduct(w http.ResponseWriter, r *http.Request)
	ListAllProduct(w http.ResponseWriter, r *http.Request)
	GetCatagoryById(w http.ResponseWriter, r *http.Request)
	ListAllBrand(w http.ResponseWriter, r *http.Request)
	UpdateCatagoryById(w http.ResponseWriter, r *http.Request)
	UpdateBrandById(w http.ResponseWriter, r *http.Request)
}

type ProductControllerImpl struct {
	productService service.ProductService
}

func NewProductController(productService service.ProductService) ProductController {
	return &ProductControllerImpl{
		productService: productService,
	}
}

func (c *ProductControllerImpl) CreateProduct(w http.ResponseWriter, r *http.Request) {
	resp, err := c.productService.CreateProduct(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to create product")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, resp)
}

func (c *ProductControllerImpl) ListAllProduct(w http.ResponseWriter, r *http.Request) {
	resp, err := c.productService.ListAllProduct(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to list all product")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, resp)
}

func (c *ProductControllerImpl) GetCatagoryById(w http.ResponseWriter, r *http.Request) {
	resp, err := c.productService.GetCatagoryById(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to get item by catagory ID")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, resp)
}

func (c *ProductControllerImpl) ListAllBrand(w http.ResponseWriter, r *http.Request) {
	resp, err := c.productService.ListAllBrands(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to get brand details")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, resp)
}

func (c *ProductControllerImpl) UpdateCatagoryById(w http.ResponseWriter, r *http.Request) {
	resp, err := c.productService.ListAllBrands(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to get brand details")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, resp)
}

func (c *ProductControllerImpl) UpdateBrandById(w http.ResponseWriter, r *http.Request) {
	resp, err := c.productService.ListAllBrands(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to get brand details")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, resp)
}
