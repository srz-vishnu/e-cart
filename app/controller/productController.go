package controller

import (
	"e-cart/app/service"
	"e-cart/pkg/api"
	"e-cart/pkg/e"
	"net/http"
)

type ProductController interface {
	CreateProduct(w http.ResponseWriter, r *http.Request)
}

type ProductControllerStructImpl struct {
	productService service.ProductService
}

func NewProductController(productService service.ProductService) ProductController {
	return ProductControllerStructImpl{
		productService: productService,
	}
}

func (c ProductControllerStructImpl) CreateProduct(w http.ResponseWriter, r *http.Request) {
	resp, err := c.productService.CreateProduct(r)
	if err != nil {
		apiErr := e.NewAPIError(err, "failed to create user")
		api.Fail(w, apiErr.StatusCode, apiErr.Code, apiErr.Message, err.Error())
		return
	}
	api.Success(w, http.StatusOK, resp)
}
