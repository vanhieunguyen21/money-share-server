package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"money_share/pkg/auth"
	"money_share/pkg/dto"
	"money_share/pkg/http/request"
	"money_share/pkg/http/response"
	"money_share/pkg/repository"
	"net/http"
	"strconv"
)

var UserRepository repository.UserRepository

func Login(w http.ResponseWriter, r *http.Request) {
	// Parse login request from body
	loginRequest := &request.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(loginRequest); err != nil {
		errMsg := fmt.Sprintf("Cannot parse request body")
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	// Validate fields
	username := loginRequest.Username
	password := loginRequest.Password
	if len(username) == 0 || len(password) == 0 {
		errMsg := fmt.Sprintf("Username or password cannot be empty")
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// Find database record and compare password
	user, err := UserRepository.GetByUsername(username)
	if err != nil {
		http.Error(w, "Wrong username or password", http.StatusUnauthorized)
		fmt.Println(err)
		return
	}
	authorized := user.ComparePassword(password)
	if !authorized {
		http.Error(w, "Wrong username or password", http.StatusUnauthorized)
		return
	}

	// Generate jwt token
	tokenStr, err := auth.GenerateJWT(username)
	if err != nil {
		errMsg := "Error when generating authorization token"
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg + " " + err.Error())
		return
	}

	// Write to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
	if err != nil {
		errMsg := fmt.Sprintf("Error encoding to json: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
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

	// Get user from database
	user, err := UserRepository.GetById(uint(userID))
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get user by ID '%d': %s", userID, err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}

	// Write to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	userDTO := dto.UserToUserDTO(*user)
	err = json.NewEncoder(w).Encode(userDTO)
	if err != nil {
		errMsg := fmt.Sprintf("Error encoding to json: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}
}

func CheckUsername(w http.ResponseWriter, r *http.Request) {
	// Init responseObj object
	responseObj := response.CheckUsernameResponse{}

	// Get username from parameters
	params := mux.Vars(r)
	username := params["username"]
	responseObj.Username = username

	// Check username requirements
	if len(username) < 6 {
		responseObj.Requirement = false
		responseObj.Message = "Length must be equal or greater than 6"
	} else {
		responseObj.Requirement = true
	}

	// Check username availability if username requirement passes
	if responseObj.Requirement {
		available, err := UserRepository.CheckUsernameAvailability(username)
		if err != nil {
			errMsg := fmt.Sprintf("Error checking username '%s' availability: %s", username, err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			fmt.Println(errMsg)
			return
		} else {
			responseObj.Available = available
			if !available {
				responseObj.Message = "Username is taken"
			}
		}
	}

	// Write to responseObj
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(responseObj)
	if err != nil {
		errMsg := fmt.Sprintf("Error encoding to json: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	// Parse user data from request body
	userDTO := &dto.UserDTO{}
	err := json.NewDecoder(r.Body).Decode(userDTO)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse request body")
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg + err.Error())
		return
	}
	user, err := userDTO.MapToDomain()
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse model: %s", err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Print(errMsg)
		return
	}

	// Validate fields
	if user.Username == "" {
		http.Error(w, "Field does not meet requirement: username must not be empty", http.StatusBadRequest)
		return
	}
	if user.Username == "" || len(user.Password) < 6 || user.DisplayName == "" {
		http.Error(w,
			"Field does not meet requirement: password length must be equal or greater than 6",
			http.StatusBadRequest)
		return
	}

	// Hash password
	err = user.HashPassword()
	if err != nil {
		errMsg := fmt.Sprintf("Error while registering")
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg + err.Error())
		return
	}

	// Create user in database
	savedUser, err := UserRepository.Create(&user)
	if err != nil {
		errMsg := fmt.Sprintf("Error creating user: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}

	// Write created user to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	savedUserDTO := dto.UserToUserDTO(*savedUser)
	err = json.NewEncoder(w).Encode(savedUserDTO)
	if err != nil {
		errMsg := fmt.Sprintf("Error encoding to json: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}
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
