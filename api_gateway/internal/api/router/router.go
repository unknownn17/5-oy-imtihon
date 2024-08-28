package router

import (
	"api/internal/config"
	"api/internal/connections"
	_ "api/internal/docs"
	jwttoken "api/internal/jwt"
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Booking Hotel API
// @version         2.0
// @description     This is an API for booking Hotels.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:8085
// @BasePath /

func NewRouter() {
	c := config.Configuration()
	r := http.NewServeMux()
	handler := connections.NewHandler()
	rate := connections.NewRateLimiting()

	// Users

	r.HandleFunc("POST /users/register", handler.Register)
	r.HandleFunc("POST /users/verify", handler.Verify)
	r.HandleFunc("POST /users/login", handler.LogIn)
	r.HandleFunc("GET /users/{id}", jwttoken.JWTMiddleware(rate.Limit(handler.GetUser)))
	r.HandleFunc("PUT /users/{id}", jwttoken.JWTMiddleware(rate.Limit(handler.UpdateUser)))
	r.HandleFunc("DELETE /users/{id}", jwttoken.JWTMiddleware(rate.Limit(handler.DeleteUser)))
	r.HandleFunc("POST /users/logout/{id}", jwttoken.JWTMiddleware(rate.Limit(handler.LogOut)))
	r.Handle("/swagger/", httpSwagger.WrapHandler)

	// Hotel

	r.HandleFunc("POST /hotels/create", jwttoken.JWTMiddleware(rate.Limit(handler.CreateHotel)))
	r.HandleFunc("POST /hotels/rooms/create", jwttoken.JWTMiddleware(rate.Limit(handler.CreateRoom)))
	r.HandleFunc("GET /hotels", jwttoken.JWTMiddleware(rate.Limit(handler.GetHotels)))
	r.HandleFunc("GET /hotels/{id}", jwttoken.JWTMiddleware(rate.Limit(handler.GetHotel)))
	r.HandleFunc("GET /hotels/rooms/{id}", jwttoken.JWTMiddleware(rate.Limit(handler.GetRooms)))
	r.HandleFunc("GET /hotels/room", jwttoken.JWTMiddleware(rate.Limit(handler.GetRoom)))
	r.HandleFunc("PUT /hotels/{id}", jwttoken.JWTMiddleware(rate.Limit(handler.UpdateHotel)))
	r.HandleFunc("PUT /hotels/rooms/{id}", jwttoken.JWTMiddleware(rate.Limit(handler.UpdateRoom)))
	r.HandleFunc("DELETE /hotels/{id}", jwttoken.JWTMiddleware(rate.Limit(handler.DeleteHotel)))
	r.HandleFunc("DELETE /hotels/rooms/{id}", jwttoken.JWTMiddleware(rate.Limit(handler.DeleteRoom)))

	//Booking

	r.HandleFunc("POST /bookings", jwttoken.JWTMiddleware(rate.Limit(handler.CreateBooking)))
	r.HandleFunc("POST /waitinglists", jwttoken.JWTMiddleware(rate.Limit(handler.CreateWaiting)))
	r.HandleFunc("GET /bookings/{id}", jwttoken.JWTMiddleware(rate.Limit(handler.GetBooking)))
	r.HandleFunc("GET /waitinglists/{id}", jwttoken.JWTMiddleware(rate.Limit(handler.GetWaiting)))
	r.HandleFunc("PUT /bookings/{id}", jwttoken.JWTMiddleware(rate.Limit(handler.UpdateBooking)))
	r.HandleFunc("PUT /waitinglists/{id}", jwttoken.JWTMiddleware(rate.Limit(handler.UpdateWaiting)))
	r.HandleFunc("DELETE /bookings/{id}", jwttoken.JWTMiddleware(rate.Limit(handler.DeleteBooking)))
	r.HandleFunc("DELETE /waitinglists/{id}", jwttoken.JWTMiddleware(rate.Limit(handler.DeleteWaiting)))

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
	err.Error()
	os.Exit(1)
}

func GracefulShutdown(srv *http.Server, logger *slog.Logger) {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, os.Kill)

	<-shutdownCh
	logger.Info("Shutdown signal received, initiating graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown error: " + err.Error())
	} else {
		logger.Info("Server gracefully stopped")
	}
}
