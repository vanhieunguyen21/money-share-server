package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"money_share/pkg/auth"
	"money_share/pkg/dto"
	"money_share/pkg/dto/request"
	"money_share/pkg/dto/response"
	"money_share/pkg/model"
	"money_share/pkg/repository"
	"net/http"
	"strconv"
)

var UserRepository repository.UserRepository

func Login(w http.ResponseWriter, r *http.Request) {
	// Parse login request from body
	loginRequest := &request.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(loginRequest); err != nil {
		ResponseError(w, "Cannot parse request body", http.StatusBadRequest)
		return
	}
	// Validate fields
	username := loginRequest.Username
	password := loginRequest.Password
	if err := model.ValidateUsername(username); err != nil {
		ResponseError(w, err.Error(), http.StatusBadRequest)
	}
	if err := model.ValidatePassword(password); err != nil {
		ResponseError(w, err.Error(), http.StatusBadRequest)
	}

	// Find database record and compare password
	user, err := UserRepository.GetByUsername(username)
	if err != nil {
		ResponseError(w, "Wrong username or password", http.StatusUnauthorized)
		return
	}
	authorized := user.ComparePassword(password)
	if !authorized {
		ResponseError(w, "Wrong username or password", http.StatusUnauthorized)
		return
	}

	// Generate jwt token
	tokenStr, err := auth.GenerateJWT(username)
	if err != nil {
		ResponseError(w, "Error when generating authorization token", http.StatusInternalServerError)
		return
	}

	// Write to response
	loginResponse := response.LoginResponse{
		UserDTO: dto.UserToUserDTO(*user),
		Token:   tokenStr,
	}
	ResponseJSON(w, loginResponse)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Get user id from parameters
	params := mux.Vars(r)
	userIDStr := params["userId"]
	userID, err := strconv.ParseUint(userIDStr, 0, 32)
	if err != nil {
		ResponseError(w, fmt.Sprintf("Cannot parse user ID '%s': %s", userIDStr, err), http.StatusBadRequest)
		return
	}

	// Get user from database
	user, err := UserRepository.GetById(uint(userID))
	if err != nil {
		ResponseError(w, fmt.Sprintf("Failed to get user by ID '%d': %s", userID, err), http.StatusInternalServerError)
		return
	}

	// Write to response
	userDTO := dto.UserToUserDTO(*user)
	ResponseJSON(w, userDTO)
}

func CheckUsername(w http.ResponseWriter, r *http.Request) {
	// Init responseObj object
	responseObj := response.SimpleResponse{}

	// Get username from parameters
	params := mux.Vars(r)
	username := params["username"]

	// Check username requirements
	if len(username) < 6 {
		responseObj.Result = false
	} else {
		available, err := UserRepository.CheckUsernameAvailability(username)
		if err != nil {
			ResponseError(w, fmt.Sprintf("Error checking username '%s' availability: %s", username, err), http.StatusInternalServerError)
			return
		} else {
			responseObj.Result = available
		}
	}

	// Write to responseObj
	ResponseJSON(w, responseObj)
}

func Register(w http.ResponseWriter, r *http.Request) {
	// Parse user data from request body
	registerRequest := &request.RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(registerRequest)
	if err != nil {
		ResponseError(w, "Cannot parse request body", http.StatusBadRequest)
		return
	}
	// Create user object
	user, err := registerRequest.UserDTO.MapToDomain()
	if err != nil {
		ResponseError(w, "Error while parsing user object", http.StatusInternalServerError)
		return
	}

	// Trim display name
	user.TrimDisplayName()

	// Validate fields
	if err = user.ValidateFields(); err != nil {
		ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set password for user
	user.Password = registerRequest.Password
	// Hash password
	if err = user.HashPassword(); err != nil {
		ResponseError(w, "Error while hashing password", http.StatusInternalServerError)
		return
	}

	// Create user in database
	savedUser, err := UserRepository.Create(&user)
	if err != nil {
		ResponseError(w, "Error while creating user", http.StatusInternalServerError)
		return
	}

	// Write created user to response
	savedUserDTO := dto.UserToUserDTO(*savedUser)
	ResponseJSON(w, savedUserDTO)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get user id from parameters
	params := mux.Vars(r)
	userIDStr := params["userId"]
	userID, err := strconv.ParseUint(userIDStr, 0, 32)
	if err != nil {
		ResponseError(w, fmt.Sprintf("Cannot parse user ID '%s': %s", userIDStr, err), http.StatusBadRequest)
		return
	}

	// Get username from header
	username := r.Header.Get("username")
	// Validate username and user id
	validated, err := UserRepository.ValidateUsernameAndUserID(username, uint(userID))
	if err != nil {
		ResponseError(w, fmt.Sprintf("Cannot validate username and user id"), http.StatusInternalServerError)
		return
	}
	if !validated {
		ResponseError(w, fmt.Sprintf("You're not authorized to do this action"), http.StatusForbidden)
		return
	}

	// Parse user data from request body
	userDTO := &dto.UserDTO{}
	if err = json.NewDecoder(r.Body).Decode(userDTO); err != nil {
		ResponseError(w, fmt.Sprintf("Cannot parse request body: %s", err), http.StatusBadRequest)
		return
	}
	userDTO.ID = uint(userID)
	user, err := userDTO.MapToDomain()
	if err != nil {
		ResponseError(w, fmt.Sprintf("Cannot parse model: %s", err), http.StatusBadRequest)
		return
	}
	// Trim display name if included
	if user.DisplayName != "" {
		user.TrimDisplayName()
	}

	// Validate all non null fields
	if err := user.ValidateNonNullFields(); err != nil {
		ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If password is included, hash it
	if user.Password != "" {
		if err := user.HashPassword(); err != nil {
			ResponseError(w, "Error while hashing password", http.StatusInternalServerError)
			return
		}
	}

	// Update user to database
	updatedUser, err := UserRepository.Update(&user)
	if err != nil {
		ResponseError(w, fmt.Sprintf("Error while updating user: %s", err), http.StatusBadRequest)
		return
	}

	// Write updated data to response
	updatedUserDTO := dto.UserToUserDTO(*updatedUser)
	ResponseJSON(w, updatedUserDTO)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get user id from parameters
	params := mux.Vars(r)
	userIDStr := params["userId"]
	userID, err := strconv.ParseUint(userIDStr, 0, 32)
	if err != nil {
		ResponseError(w, fmt.Sprintf("Cannot parse user ID '%s': %s", userIDStr, err), http.StatusBadRequest)
		return
	}

	// Get username from header
	username := r.Header.Get("username")
	// Validate username and user id
	validated, err := UserRepository.ValidateUsernameAndUserID(username, uint(userID))
	if err != nil {
		ResponseError(w, "Cannot validate username and user id", http.StatusInternalServerError)
		return
	}
	if !validated {
		ResponseError(w, "You're not authorized to do this action", http.StatusForbidden)
		return
	}

	// Delete user from database and write response
	if err := UserRepository.Delete(uint(userID)); err != nil {
		ResponseError(w,  fmt.Sprintf("Error while deleting user with ID '%d'", userID), http.StatusInternalServerError)
		return
	}

	// Write to response
	responseObj := response.SimpleResponse{Result: true}
	ResponseJSON(w, responseObj)
}
