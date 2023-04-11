// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"

	"github.com/RipperAcskt/innotaxidriver/config"
	user "github.com/RipperAcskt/innotaxidriver/internal/client"
	"github.com/RipperAcskt/innotaxidriver/internal/handler/grpc"
	handler "github.com/RipperAcskt/innotaxidriver/internal/handler/restapi"
	"github.com/RipperAcskt/innotaxidriver/internal/repo/cassandra"
	"github.com/RipperAcskt/innotaxidriver/internal/service"
	"github.com/RipperAcskt/innotaxidriver/restapi/operations"
	"github.com/RipperAcskt/innotaxidriver/restapi/operations/auth"
	"github.com/RipperAcskt/innotaxidriver/restapi/operations/driver"
)

//go:generate swagger generate server --target ../../InnoTaxiDriver --name InnoTaxiDriverAPI --spec ../docs/swagger.yaml --principal interface{} --default-scheme gin-swwagger

func configureFlags(api *operations.InnoTaxiDriverAPIAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.InnoTaxiDriverAPIAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError
	api.UseSwaggerUI()
	log, err := zap.NewProduction()
	if err != nil {
		log.Sugar().Fatalf("config new failed: %v", err)
	}

	cfg, err := config.New()
	if err != nil {
		log.Sugar().Fatalf("config new failed: %v", err)
	}

	cassandra, err := cassandra.New(cfg)
	if err != nil {
		log.Sugar().Fatalf("cassandra new failed: %v", err)
	}

	err = cassandra.M.Up()
	if err != migrate.ErrNoChange && err != nil {
		log.Sugar().Fatalf("migrate up failed: %v", err)
	}

	client, err := user.New(cfg)
	if err != nil {
		log.Sugar().Fatalf("grpc new failed: %v", err)
	}

	service := service.New(cassandra, client, cfg)
	handler := handler.New(service, cfg)

	api.AuthPostDriverSingUpHandler = auth.PostDriverSingUpHandlerFunc(handler.SingUp)
	api.AuthPostDriverSingInHandler = auth.PostDriverSingInHandlerFunc(handler.SingIn)
	api.AuthPostDriverRefreshHandler = auth.PostDriverRefreshHandlerFunc(handler.Refresh)

	api.DriverGetDriverHandler = driver.GetDriverHandlerFunc(handler.GetProfile)
	api.DriverPutDriverHandler = driver.PutDriverHandlerFunc(handler.UpdateProfile)
	api.DriverDeleteDriverHandler = driver.DeleteDriverHandlerFunc(handler.DeleteProfile)

	api.AddMiddlewareFor("POST", "/driver/refresh", handler.VerifyToken)
	api.AddMiddlewareFor("GET", "/driver", handler.VerifyToken)
	api.AddMiddlewareFor("PUT", "/driver", handler.VerifyToken)
	api.AddMiddlewareFor("DELETE", "/driver", handler.VerifyToken)

	api.AddMiddlewareFor("POST", "/driver", handler.Recovery)
	api.AddMiddlewareFor("GET", "/driver", handler.Recovery)
	api.AddMiddlewareFor("PUT", "/driver", handler.Recovery)
	api.AddMiddlewareFor("DELETE", "/driver", handler.Recovery)

	api.AddMiddlewareFor("POST", "/driver/sing-up", handler.Log)
	api.AddMiddlewareFor("POST", "/driver/sing-in", handler.Log)
	api.AddMiddlewareFor("POST", "/driver/refresh", handler.Log)
	api.AddMiddlewareFor("GET", "/driver", handler.Log)
	api.AddMiddlewareFor("PUT", "/driver", handler.Log)
	api.AddMiddlewareFor("DELETE", "/driver", handler.Log)

	api.JSONConsumer = runtime.JSONConsumer()

	if api.AuthPostDriverSingUpHandler == nil {
		api.AuthPostDriverSingUpHandler = auth.PostDriverSingUpHandlerFunc(func(params auth.PostDriverSingUpParams) middleware.Responder {
			return middleware.NotImplemented("operation auth.PostDriverSingUp has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	grpcServer := grpc.New(service.Order, cfg)
	go func() {
		if err := grpcServer.Run(); err != nil {
			log.Error(fmt.Sprintf("grpc server run failed: %v", err))
			return
		}
	}()

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
