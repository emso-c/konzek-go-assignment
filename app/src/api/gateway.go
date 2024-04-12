// Package api provides functionality for initializing HTTP API routes and registering
// middleware handlers for handling various tasks such as CORS, CSRF protection, rate limiting,
// and SQL injection prevention.
package api

import (
	"github.com/emso-c/konzek-go-assignment/src/api/middlewares"
	"github.com/emso-c/konzek-go-assignment/src/api/routers"
	"github.com/emso-c/konzek-go-assignment/src/modules/limiter"
	"github.com/gorilla/mux"
)

var router *mux.Router

// Init initializes the HTTP API routes and registers middleware handlers.
func Init() {
	router = mux.NewRouter()

	router = router.PathPrefix("/api").Subrouter()

	// Register routers
	routers.RegisterTasksRouter(router)
	limiter.GetLimiter().Initialize()

	// TODO: Add authentication & authorization middlewares.
	router.NotFoundHandler = middlewares.NotFoundMiddleware()
	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.CSRFMiddleware())
	router.Use(middlewares.RateLimitMiddleware())
	router.Use(middlewares.SQLInjectionMiddleware())
}

// GetRouter retrieves the initialized router instance.
func GetRouter() *mux.Router {
	return router
}
