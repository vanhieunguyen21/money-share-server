package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"money_share/pkg/auth"
	"money_share/pkg/dto"
	"money_share/pkg/dto/request"
	"money_share/pkg/dto/response"
	"money_share/pkg/repository"
	"money_share/pkg/util"
	"net/http"
	"strconv"
)

var UserRepository repository.UserRepository

func Login(w http.ResponseWriter, r *http.Request) {
	// Parse login request from body
	loginRequest := &request.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(loginRequest); err != nil {
		util.ResponseError(w, "Cannot parse request body", http.StatusBadRequest)
		return
	}
	// Validate fields
	username := loginRequest.Username
	password := loginRequest.Password
	if len(username) == 0 || len(password) == 0 {
		util.ResponseError(w, "Username or password cannot be empty", http.StatusBadRequest)
		return
	}

	// Find database record and compare password
	user, err := UserRepository.GetByUsername(username)
	if err != nil {
		util.ResponseError(w, "Wrong username or password", http.StatusUnauthorized)
		return
	}
	authorized := user.ComparePassword(password)
	if !authorized {
		util.ResponseError(w, "Wrong username or password", http.StatusUnauthorized)
		return
	}

	// Generate jwt token
	tokenStr, err := auth.GenerateJWT(username)
	if err != nil {
		util.ResponseError(w, "Error when generating authorization token", http.StatusInternalServerError)
		return
	}

	// Write to response
	loginResponse := response.LoginResponse{
		UserDTO: dto.UserToUserDTO(*user),
		Token:   tokenStr,
	}
	util.ResponseJSON(w, loginResponse)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Get user id from parameters
	params := mux.Vars(r)
	userIDStr := params["userId"]
	userID, err := strconv.ParseUint(userIDStr, 0, 32)
	if err != nil {
		util.ResponseError(w, fmt.Sprintf("Cannot parse user ID '%s': %s", userIDStr, err), http.StatusBadRequest)
		return
	}

	// Get user from database
	user, err := UserRepository.GetById(uint(userID))
	if err != nil {
		util.ResponseError(w, fmt.Sprintf("Failed to get user by ID '%d': %s", userID, err), http.StatusInternalServerError)
		return
	}

	// Write to response
	userDTO := dto.UserToUserDTO(*user)
	util.ResponseJSON(w, userDTO)
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
			util.ResponseError(w, fmt.Sprintf("Error checking username '%s' availability: %s", username, err), http.StatusInternalServerError)
			return
		} else {
			responseObj.Result = available
		}
	}

	// Write to responseObj
	util.ResponseJSON(w, responseObj)
}

func Register(w http.ResponseWriter, r *http.Request) {
	// Parse user data from request body
	registerRequest := &request.RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(registerRequest)
	if err != nil {
		util.ResponseError(w, "Cannot parse request body", http.StatusBadRequest)
		return
	}

	// Validate fields
	if len(registerRequest.Username) < 6 {
		util.ResponseError(w, "Field does not meet requirement: username must be at least 6 characters", http.StatusBadRequest)
		return
	}
	if len(registerRequest.Password) < 8 {
		util.ResponseError(w, "Field does not meet requirement: password must be at least 8 characters", http.StatusBadRequest)
		return
	}
	if len(registerRequest.DisplayName) < 4 {
		util.ResponseError(w, "Field does not meet requirement: display name must be at least 4 characters", http.StatusBadRequest)
		return
	}

	// Create user object
	user, err := registerRequest.UserDTO.MapToDomain()
	if err != nil {
		util.ResponseError(w, "Error while parsing user object", http.StatusInternalServerError)
		return
	}
	user.Password = registerRequest.Password

	// Hash password
	err = user.HashPassword()
	if err != nil {
		util.ResponseError(w, "Error while registering", http.StatusInternalServerError)
		return
	}

	// Create user in database
	savedUser, err := UserRepository.Create(&user)
	if err != nil {
		util.ResponseError(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Write created user to response
	savedUserDTO := dto.UserToUserDTO(*savedUser)
	util.ResponseJSON(w, savedUserDTO)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get user id from parameters
	params := mux.Vars(r)
	userIDStr := params["userId"]
	userID, err := strconv.ParseUint(userIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse user ID '%s': %s", userIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Get username from header
	username := r.Header.Get("username")
	// Validate username and user id
	validated, err := UserRepository.ValidateUsernameAndUserID(username, uint(userID))
	if err != nil {
		errMsg := fmt.Sprintf("Cannot validate username and user id")
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	if !validated {
		errMsg := fmt.Sprintf("You're not authorized to do this action")
		http.Error(w, errMsg, http.StatusForbidden)
		return
	}

	// Parse user data from request body
	userDTO := &dto.UserDTO{}
	err = json.NewDecoder(r.Body).Decode(userDTO)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse request body: %s", err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}
	userDTO.ID = uint(userID)
	user, err := userDTO.MapToDomain()
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse model: %s", err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Print(errMsg)
		return
	}

	// Update user to database
	updatedUser, err := UserRepository.Update(&user)
	if err != nil {
		errMsg := fmt.Sprintf("Error while updating user: %s", err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Write updated data to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	updatedUserDTO := dto.UserToUserDTO(*updatedUser)
	err = json.NewEncoder(w).Encode(updatedUserDTO)
	if err != nil {
		errMsg := fmt.Sprintf("Error encoding to json: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get user id from parameters
	params := mux.Vars(r)
	userIDStr := params["userId"]
	userID, err := strconv.ParseUint(userIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse user ID '%s': %s", userIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Get username from header
	username := r.Header.Get("username")
	// Validate username and user id
	validated, err := UserRepository.ValidateUsernameAndUserID(username, uint(userID))
	if err != nil {
		errMsg := fmt.Sprintf("Cannot validate username and user id")
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	if !validated {
		errMsg := fmt.Sprintf("You're not authorized to do this action")
		http.Error(w, errMsg, http.StatusForbidden)
		return
	}

	// Delete user from database and write response
	err = UserRepository.Delete(uint(userID))
	if err != nil {
		errMsg := fmt.Sprintf("Error deleting user with ID '%d': %s", userID, err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}

	// Write to response
	w.WriteHeader(http.StatusOK)
}
