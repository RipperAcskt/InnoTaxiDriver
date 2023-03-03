// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/golang-migrate/migrate/v4"

	"github.com/RipperAcskt/innotaxidriver/config"
	user "github.com/RipperAcskt/innotaxidriver/internal/client"
	"github.com/RipperAcskt/innotaxidriver/internal/handler"
	"github.com/RipperAcskt/innotaxidriver/internal/repo/cassandra"
	"github.com/RipperAcskt/innotaxidriver/internal/service"
	"github.com/RipperAcskt/innotaxidriver/restapi/operations"
	"github.com/RipperAcskt/innotaxidriver/restapi/operations/auth"
)

//go:generate swagger generate server --target ../../InnoTaxiDriver --name InnoTaxiDriverAPI --spec ../docs/swagger.yaml --principal interface{} --default-scheme gin-swwagger

func configureFlags(api *operations.InnoTaxiDriverAPIAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.InnoTaxiDriverAPIAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError
	api.UseSwaggerUI()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config new failed: %v", err)
	}

	cassandra, err := cassandra.New(cfg)
	if err != nil {
		log.Fatalf("cassandra new failed: %v", err)
	}

	err = cassandra.M.Up()
	if err != migrate.ErrNoChange && err != nil {
		log.Fatalf("migrate up failed: %v", err)
	}

	client, err := user.New(cfg)
	if err != nil {
		log.Fatalf("grpc new failed: %v", err)
	}

	service := service.New(cassandra, client, cfg)
	handler := handler.New(service, cfg)

	api.AuthPostDriverSingUpHandler = auth.PostDriverSingUpHandlerFunc(handler.SingUp)
	api.AuthPostDriverSingInHandler = auth.PostDriverSingInHandlerFunc(handler.SingIn)
	api.AuthPostDriverRefreshHandler = auth.PostDriverRefreshHandlerFunc(handler.Refresh)

	api.AddMiddlewareFor("POST", "/driver/refresh", handler.VerifyToken)

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.AuthPostDriverSingUpHandler == nil {
		api.AuthPostDriverSingUpHandler = auth.PostDriverSingUpHandlerFunc(func(params auth.PostDriverSingUpParams) middleware.Responder {
			return middleware.NotImplemented("operation auth.PostDriverSingUp has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
