package routers

import (
	"github.com/product-mgmt/common-service/constants/endpoints"
	"github.com/product-mgmt/common-service/middleware"
	"github.com/product-mgmt/common-service/utils/commfunc"
	"github.com/product-mgmt/product-service/internal/handlers"
)

func (s *Storage) RegisterAuthRoutes() {
	ctrl := handlers.New(s.logger, s.sqlStore)
	midd := middleware.New(s.logger, s.sqlStore)

	publicRoute := s.router.PathPrefix(endpoints.PRODUCT_BASE_PATH).Subrouter()
	privateRoute := s.router.PathPrefix(endpoints.PRODUCT_BASE_PATH).Subrouter()

	privateRoute.Use(midd.Authenticate)
	privateRoute.Use(midd.Authorization("admin"))

	// public routes
	publicRoute.HandleFunc(endpoints.LIST, commfunc.MakeHTTPHandleFunc(ctrl.GetProductsHandler)).Methods("GET")
	publicRoute.HandleFunc(endpoints.VIEW, commfunc.MakeHTTPHandleFunc(ctrl.GetProductsHandler)).Methods("GET")

	publicRoute.HandleFunc(endpoints.CATEGORY+endpoints.LIST, commfunc.MakeHTTPHandleFunc(ctrl.GetProductsHandler)).Methods("GET")
	publicRoute.HandleFunc(endpoints.CATEGORY+endpoints.VIEW, commfunc.MakeHTTPHandleFunc(ctrl.GetProductsHandler)).Methods("GET")

	publicRoute.HandleFunc(endpoints.INVENTORY+endpoints.LIST, commfunc.MakeHTTPHandleFunc(ctrl.GetProductsHandler)).Methods("GET")
	publicRoute.HandleFunc(endpoints.INVENTORY+endpoints.VIEW, commfunc.MakeHTTPHandleFunc(ctrl.GetProductsHandler)).Methods("GET")

	// // private routes
	privateRoute.HandleFunc(endpoints.ADD, commfunc.MakeHTTPHandleFunc(ctrl.GetProductsHandler)).Methods("GET")
	privateRoute.HandleFunc(endpoints.UPDATE, commfunc.MakeHTTPHandleFunc(ctrl.GetProductsHandler)).Methods("GET")
	privateRoute.HandleFunc(endpoints.DELETE, commfunc.MakeHTTPHandleFunc(ctrl.GetProductsHandler)).Methods("GET")

	privateRoute.HandleFunc(endpoints.CATEGORY+endpoints.ADD, commfunc.MakeHTTPHandleFunc(ctrl.GetProductsHandler)).Methods("GET")
	privateRoute.HandleFunc(endpoints.CATEGORY+endpoints.UPDATE, commfunc.MakeHTTPHandleFunc(ctrl.GetProductsHandler)).Methods("GET")
	privateRoute.HandleFunc(endpoints.CATEGORY+endpoints.DELETE, commfunc.MakeHTTPHandleFunc(ctrl.GetProductsHandler)).Methods("GET")

	privateRoute.HandleFunc(endpoints.INVENTORY+endpoints.ADD, commfunc.MakeHTTPHandleFunc(ctrl.GetProductsHandler)).Methods("GET")
	privateRoute.HandleFunc(endpoints.INVENTORY+endpoints.UPDATE, commfunc.MakeHTTPHandleFunc(ctrl.GetProductsHandler)).Methods("GET")
	privateRoute.HandleFunc(endpoints.INVENTORY+endpoints.DELETE, commfunc.MakeHTTPHandleFunc(ctrl.GetProductsHandler)).Methods("GET")
}
