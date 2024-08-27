package router

import (
	"api/internal/config"
	"api/internal/connections"
	_ "api/internal/docs"
	jwttoken "api/internal/jwt"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           ITEMS API
// @version         2.0
// @description     This is an API for booking Hotels.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:8085
// @BasePath /


func NewRouter() {
	c := config.Configuration()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	r := http.NewServeMux()
	handler := connections.NewHandler()

	r.HandleFunc("POST /users/register", handler.Register)
	r.HandleFunc("POST /users/verify", handler.Verify)
	r.HandleFunc("POST /users/login", handler.LogIn)
	r.HandleFunc("GET /users/{id}", jwttoken.JWTMiddleware(handler.GetUser))
	r.HandleFunc("PUT /users/{id}", jwttoken.JWTMiddleware(handler.UpdateUser))
	r.HandleFunc("DELETE /users/{id}", jwttoken.JWTMiddleware(handler.DeleteUser))
	r.HandleFunc("POST /users/logout/{id}", jwttoken.JWTMiddleware(handler.LogOut))
	r.Handle("/swagger/", httpSwagger.WrapHandler)

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	srv := &http.Server{
		Addr:      c.User.Port,
		Handler:   r,
		TLSConfig: tlsConfig,
	}
	fmt.Printf("Server started on port %s\n", c.User.Port)
	err := srv.ListenAndServeTLS("./tls/localhost.pem", "./tls/localhost-key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}
