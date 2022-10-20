package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"money_share/pkg/auth"
	"money_share/pkg/dto"
	"money_share/pkg/dto/request"
	"money_share/pkg/dto/response"
	"money_share/pkg/model"
	"money_share/pkg/repository"
	"money_share/pkg/util"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 2 MB
const maxUploadSize = 2 << 20

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
		return
	}
	if err := model.ValidatePassword(password); err != nil {
		ResponseError(w, err.Error(), http.StatusBadRequest)
		return
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
	accessToken, refreshToken, err := auth.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		ResponseError(w, "Error when generating tokens", http.StatusInternalServerError)
		return
	}

	// Write to response
	loginResponse := response.LoginResponse{
		UserDTO:      dto.UserToUserDTO(*user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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
	if len(username) < 8 {
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
	// Set password for user
	user.Password = registerRequest.Password

	// Trim display name
	user.TrimDisplayName()

	// Validate fields
	if err = user.ValidateFields(); err != nil {
		ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password
	if err = user.HashPassword(); err != nil {
		ResponseError(w, "Error while hashing password", http.StatusInternalServerError)
		return
	}

	// Create user in database
	err = UserRepository.Create(&user)
	if err != nil {
		ResponseError(w, "Error while creating user", http.StatusInternalServerError)
		return
	}

	// Write created user to response
	responseObj := response.SimpleResponse{Result: true}
	ResponseJSON(w, responseObj)
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
		ResponseError(w, fmt.Sprintf("You don't have permission to do this action"), http.StatusForbidden)
		return
	}

	// Parse user data from request body
	updateUserRequest := &request.UpdateUserRequest{}
	if err = json.NewDecoder(r.Body).Decode(updateUserRequest); err != nil {
		ResponseError(w, fmt.Sprintf("Cannot parse request body: %s", err), http.StatusBadRequest)
		return
	}
	updateMap := make(map[string]interface{})
	// Parse fields
	// Display name
	if updateUserRequest.DisplayName != nil {
		displayName := strings.TrimSpace(*updateUserRequest.DisplayName)
		if err := model.ValidateDisplayName(displayName); err != nil {
			ResponseError(w, err.Error(), http.StatusBadRequest)
			return
		}
		updateMap["DisplayName"] = displayName
	}
	// Password
	if updateUserRequest.Password != nil {
		if err := model.ValidatePassword(*updateUserRequest.Password); err != nil {
			ResponseError(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Hash password
		hashedPassword, err := model.HashPassword(*updateUserRequest.Password)
		if err != nil {
			ResponseError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		updateMap["Password"] = hashedPassword
	}
	// Phone number
	if updateUserRequest.PhoneNumber != nil {
		updateMap["PhoneNumber"] = *updateUserRequest.PhoneNumber
	}
	// Email address
	if updateUserRequest.EmailAddress != nil {
		updateMap["EmailAddress"] = *updateUserRequest.EmailAddress
	}
	// Date of birth
	if updateUserRequest.DateOfBirth != nil {
		dob, err := time.Parse(util.ShortDateLayout, *updateUserRequest.DateOfBirth)
		if err != nil {
			ResponseError(w, fmt.Sprintf("Cannot parse date of birth: %s", err), http.StatusBadRequest)
			return
		}
		updateMap["DateOfBirth"] = dob
	}

	// Update user to database
	updatedUser, err := UserRepository.Update(uint(userID), updateMap)
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
		ResponseError(w, fmt.Sprintf("Error while deleting user with ID '%d'", userID), http.StatusInternalServerError)
		return
	}

	// Write to response
	responseObj := response.SimpleResponse{Result: true}
	ResponseJSON(w, responseObj)
}

func UploadUserProfileImage(w http.ResponseWriter, r *http.Request) {
	// Limit upload file size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		ResponseError(w, "Uploaded file too big. Max file size is 2MB.", http.StatusBadRequest)
		return
	}

	// Get file from request
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		ResponseError(w, fmt.Sprintf("Error getting file from request: %s", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

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

	// Create upload folder if it doesn't exist
	err = os.MkdirAll("./fileServer/userProfileImage", os.ModePerm)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create new file in upload folder
	profileImageUrl := fmt.Sprintf("./fileServer/userProfileImage/%s_%d%s",
		username, time.Now().Unix(), filepath.Ext(fileHeader.Filename))
	dst, err := os.Create(profileImageUrl)
	defer dst.Close()

	// Copy uploaded file to create file
	_, err = io.Copy(dst, file)
	if err != nil {
		ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update profile image url for user
	updateMap := make(map[string]interface{})
	updateMap["ProfileImageUrl"] = filepath.Base(profileImageUrl)
	updatedUser, err := UserRepository.Update(uint(userID), updateMap)
	if err != nil {
		ResponseError(w, "Error while updating profile image", http.StatusInternalServerError)
		return
	}

	// Write updated data to response
	updatedUserDTO := dto.UserToUserDTO(*updatedUser)
	ResponseJSON(w, updatedUserDTO)
}

func GetUserProfileImage(w http.ResponseWriter, r *http.Request) {
	// Get file name from parameters
	params := mux.Vars(r)
	fileName := params["fileName"]

	fullFilePath := "./fileServer/userProfileImage/" + fileName
	// Check if file exists
	if _, err := os.Stat(fullFilePath); err != nil {
		ResponseError(w, "File doesn't exist", http.StatusBadRequest)
		return
	}

	// Read file
	fileBytes, err := os.ReadFile(fullFilePath)
	if err != nil {
		ResponseError(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Write to response
	ResponseFile(w, fileBytes)
}
