package app

import (
	"e-cart/app/controller"
	"e-cart/app/internal"
	"e-cart/app/service"
	api "e-cart/pkg/api"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func APIRouter(db *gorm.DB) chi.Router {
	r := chi.NewRouter()

	// User part
	urRepo := internal.NewUserRepo(db)
	urService := service.NewUserService(urRepo)
	urController := controller.NewUserController(urService)

	// Product part
	proRepo := internal.NewProductRepo(db)
	proService := service.NewProductService(proRepo)
	proController := controller.NewProductController(proService)

	// Admin part
	adminRepo := internal.NewAdminRepo(db)
	adminService := service.NewAdminService(adminRepo)
	adminController := controller.NewAdminController(adminService)

	r.Route("/", func(r chi.Router) {
		r.Get("/hello", api.ExampleHamdler)
	})

	//user
	r.Route("/user", func(r chi.Router) {
		r.Post("/create", urController.UserDetails)
		r.Post("/login", urController.LoginUser)
		r.Put("/update/{userid}", urController.UpdateUserDetails)
		r.Post("/cart/additem", urController.AddItemsToCart)
		r.Post("/cart/placeorder", urController.PlaceOrder)
	})

	//product
	r.Route("/product", func(r chi.Router) {
		// create used to create a product
		r.Post("/create", proController.CreateProduct)
		// to view product
		r.Get("/list/catagory", proController.ListAllProduct)
		// to view brand
		r.Get("/list/brand", proController.ListAllBrand)
		// to see product based on the catagory id given by user
		r.Get("/search/catagory/{id}", proController.GetCatagoryById)

		// r.Put("/list/catagory/{id}", proController.UpdateBrandById)
		// r.Put("/list/brand/{id}", proController.UpdateBrand)

	})

	//admin
	r.Route("/admin", func(r chi.Router) {
		r.Put("/block/{useridid}", adminController.BlockUser)
		r.Put("/unblock/{useridid}", adminController.BlockUser)
	})

	return r
}
