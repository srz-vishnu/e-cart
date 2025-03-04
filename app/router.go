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

	r.Route("/", func(r chi.Router) {
		r.Get("/hello", api.ExampleHamdler)
	})

	//user
	r.Route("/user", func(r chi.Router) {
		r.Post("/create", urController.UserDetails)
	})

	//product
	r.Route("/product", func(r chi.Router) {
		r.Post("/create", proController.CreateProduct)

	})

	return r
}
