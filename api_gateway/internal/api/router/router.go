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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	r := http.NewServeMux()
	handler := connections.NewHandler()

	// Users

	r.HandleFunc("POST /users/register", handler.Register)
	r.HandleFunc("POST /users/verify", handler.Verify)
	r.HandleFunc("POST /users/login", handler.LogIn)
	r.HandleFunc("GET /users/{id}", jwttoken.JWTMiddleware(handler.GetUser))
	r.HandleFunc("PUT /users/{id}", jwttoken.JWTMiddleware(handler.UpdateUser))
	r.HandleFunc("DELETE /users/{id}", jwttoken.JWTMiddleware(handler.DeleteUser))
	r.HandleFunc("POST /users/logout/{id}", jwttoken.JWTMiddleware(handler.LogOut))
	r.Handle("/swagger/", httpSwagger.WrapHandler)

	// Hotel

	r.HandleFunc("POST /hotels/create",jwttoken.JWTMiddleware(handler.CreateHotel))
	r.HandleFunc("POST /hotels/rooms/create",jwttoken.JWTMiddleware(handler.CreateRoom))
	r.HandleFunc("GET /hotels",jwttoken.JWTMiddleware(handler.GetHotels))
	r.HandleFunc("GET /hotels/{id}",jwttoken.JWTMiddleware(handler.GetHotel))
	r.HandleFunc("GET /hotels/rooms/{id}",jwttoken.JWTMiddleware(handler.GetRooms))
	r.HandleFunc("GET /hotels/room",jwttoken.JWTMiddleware(handler.GetRoom))
	r.HandleFunc("PUT /hotels/{id}",jwttoken.JWTMiddleware(handler.UpdateHotel))
	r.HandleFunc("PUT /hotels/rooms/{id}",jwttoken.JWTMiddleware(handler.UpdateRoom))
	r.HandleFunc("DELETE /hotels/{id}",jwttoken.JWTMiddleware(handler.DeleteHotel))
	r.HandleFunc("DELETE /hotels/rooms/{id}",jwttoken.JWTMiddleware(handler.DeleteRoom))

	//Booking

	r.HandleFunc("POST /bookings",jwttoken.JWTMiddleware(handler.CreateBooking))
	r.HandleFunc("POST /waitinglists",jwttoken.JWTMiddleware(handler.CreateWaiting))
	r.HandleFunc("GET /bookings/{id}",jwttoken.JWTMiddleware(handler.GetBooking))
	r.HandleFunc("GET /waitinglists/{id}",jwttoken.JWTMiddleware(handler.GetWaiting))
	r.HandleFunc("PUT /bookings/{id}",jwttoken.JWTMiddleware(handler.UpdateBooking))
	r.HandleFunc("PUT /waitinglists/{id}",jwttoken.JWTMiddleware(handler.UpdateWaiting))
	r.HandleFunc("DELETE /bookings/{id}",jwttoken.JWTMiddleware(handler.DeleteBooking))
	r.HandleFunc("DELETE /waitinglists/{id}",jwttoken.JWTMiddleware(handler.DeleteWaiting))

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

  