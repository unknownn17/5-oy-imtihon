package handler

import (
	broadcast17 "api/internal/broadcast"
	"api/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	B *broadcast17.Adjust
}

// Register handles the user registration process.
// @Summary Register a new user
// @Description Register a new user by providing email and password. A verification code will be sent to the provided email.
// @Tags user
// @Accept json
// @Produce json
// @Param registerUserRequest body models.RegisterUserRequest true "User registration data"
// @Success 200 {string} string "Verification code is sent to your email"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 403 {string} string "Already Exists"
// @Router /users/register [post]
func (u *Handler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.RegisterUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := u.B.Register(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Verification code is sent to your email")
}

// Verify handles the user verification process.
// @Summary Verify a user account
// @Description Verify a user account by providing the verification code sent to the user's email.
// @Tags user
// @Accept json
// @Produce json
// @Param verifyRequest body models.VerifyRequest true "Verification code data"
// @Success 200 {string} string "You have verified your account and now you can log in"
// @Failure 500 {string} string "Internal Server Error"
// @Router /users/verify [post]
func (u *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.VerifyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := u.B.Verify(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("You have verified your account and now you can log in")
}

// LogIn handles the user login process.
// @Summary Log in a user
// @Description Log in a user by providing their email and password.
// @Tags user
// @Accept json
// @Produce json
// @Param logInRequest body models.LogInRequest true "User login data"
// @Success 200 {string} string "JWT token"
// @Failure 500 {string} string "Internal Server Error"
// @Router /users/login [post]
func (u *Handler) LogIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.LogInRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := u.B.Login(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(token)
}

// GetUser retrieves user information by ID.
// @Summary Get user information
// @Description Retrieve user information by providing the user ID.
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.GetUserResponse "User information"
// @Failure 500 {string} string "Internal Server Error"
// @Security BearerAuth
// @Router /users/{id} [get]
func (u *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := u.B.GetUser(&models.GetUserRequest{ID: int32(id)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

// UpdateUser updates user information by ID.
// @Summary Update user information
// @Description Update user information by providing the user ID in the path and the fields to update in the request body.
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param updateUserRequest body models.UpdateUserRequest true "User update data"
// @Success 200 {string} string "Your account is updating we'll notify you when it's updated"
// @Failure 500 {string} string "Internal Server Error"
// @Security BearerAuth
// @Router /users/{id} [put]
func (u *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var req models.UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.ID = int32(id)
	if err := u.B.UpdateUser(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Your account is updating we'll notify you when it's updated")
}

// DeleteUser deletes a user by ID.
// @Summary Delete a user
// @Description Delete a user by providing their user ID.
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {string} string "Your account is deleting"
// @Failure 500 {string} string "Internal Server Error"
// @Security BearerAuth
// @Router /users/{id} [delete]
func (u *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := u.B.DeleteUser(&models.GetUserRequest{ID: int32(id)}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Your account is deleting")
}

// LogOut logs out a user by ID.
// @Summary Log out a user
// @Description Log out a user by providing their user ID.
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {string} string "You have logged out!"
// @Failure 500 {string} string "Internal Server Error"
// @Security BearerAuth
// @Router /users/logout/{id} [post]
func (u *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := u.B.Logout(&models.GetUserRequest{ID: int32(id)}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("You have logged out!")
}

// Hotel---Service

func (u *Handler) CreateHotel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.CreateHotelRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := u.B.CreateHotel(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Hotel created successfully")
}

func (u *Handler) GetHotel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := u.B.GetHotel(&models.GetHotelRequest{ID: int32(id)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (u *Handler) GetHotels(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	res,err:=u.B.GetHotels(&models.GetsRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}
