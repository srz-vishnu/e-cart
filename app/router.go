package app

import (
	"e-cart/app/controller"
	"e-cart/app/helper"
	"e-cart/app/internal"
	"e-cart/app/service"
	api "e-cart/pkg/api"
	"e-cart/pkg/middleware"
	"e-cart/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"gorm.io/gorm"
)

func APIRouter(db *gorm.DB) chi.Router {
	r := chi.NewRouter()

	// User part
	urRepo := internal.NewUserRepo(db)
	hlRepo := helper.NewContextHelper()
	hashPkg := utils.NewBcryptPackage()
	urService := service.NewUserService(urRepo, hlRepo, hashPkg)
	urController := controller.NewUserController(urService)

	// Product part
	proRepo := internal.NewProductRepo(db)
	proService := service.NewProductService(proRepo)
	proController := controller.NewProductController(proService)

	// Admin part
	adminRepo := internal.NewAdminRepo(db)
	adminService := service.NewAdminService(adminRepo, urRepo)
	adminController := controller.NewAdminController(adminService)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	r.Route("/", func(r chi.Router) {
		r.Get("/hello", api.ExampleHamdler)
		r.Post("/signup", urController.UserDetails)
		r.Post("/login", urController.LoginUser)
	})

	// User routes — JWT middleware applied
	r.Route("/user", func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware) // All user routes require login
		r.Put("/update/{userid}", urController.UpdateUserDetails)
		r.Post("/change/pwd", urController.ChangePassword)
		r.Get("/{userid}", urController.GetUserDetails)
		r.Post("/cart/additem", urController.AddItemsToCart)
		r.Get("/cart/view", urController.ViewUserCart)
		r.Delete("/cart/clear", urController.ClearCart)
		r.Post("/cart/placeorder", urController.PlaceOrder)
		r.Get("/order/history", urController.OrderHistory)
		r.Post("/favourite", urController.AddItemsToFavourites)
		r.Get("/favourite", urController.GetUserFavouriteItems)
	})

	// Product routes — JWT middleware applied
	r.Route("/product", func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware) //  All product routes need login

		r.Get("/list/catagory", proController.ListAllProduct)
		r.Get("/list/brand", proController.ListAllBrand)
		r.Get("/category/{id}", proController.GetBrandByID)
		r.Get("/search/catagory/id/{id}", proController.GetCatagoryById)
		r.Get("/catagory/id/{id}", proController.GetCatagoryDetailsById)
		r.Get("/search/catagory/name/{categoryname}", proController.GetCatagoryByName)

		// Create product — admin only
		r.With(middleware.AdminOnlyMiddleware).Post("/create", proController.CreateProduct)
	})

	// Admin routes — JWT and Admin middleware
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware)   //  Require login
		r.Use(middleware.AdminOnlyMiddleware) //  Must be admin

		r.Put("/block/{userid}", adminController.BlockUser)     // admin only
		r.Put("/unblock/{userid}", adminController.UnBlockUser) // admin only
		r.Get("/userdetails", adminController.GetAllUserDetail)
		r.Get("/block/userdetails", adminController.GetAllBlockedUserDetail) //admin only
		r.Get("/order/history/{id}", adminController.CustomerOrderHistoryById)
		r.Get("/getall/order/history", adminController.CustomerOrderHistory)
	})

	return r
}
