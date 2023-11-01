package routers

import (
	"github.com/ankeshnirala/order-mgmt/product-service/internal/handlers"
	"github.com/product-mgmt/common-service/constants/endpoints"
	"github.com/product-mgmt/common-service/middleware"
	"github.com/product-mgmt/common-service/utils/commfunc"
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
	// publicRoute.HandleFunc(endpoints.SIGNIN_PATH, commfunc.MakeHTTPHandleFunc(ctrl.SigninHandler)).Methods("POST")

	// // private routes
	// privateRoute.HandleFunc(endpoints.PROFILE_PATH, commfunc.MakeHTTPHandleFunc(ctrl.ProfileHandler)).Methods("GET")
}
