package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"money_share/pkg/dto"
	"money_share/pkg/dto/request"
	"money_share/pkg/model"
	"money_share/pkg/repository"
	"net/http"
	"strconv"
)

var GroupRepository repository.GroupRepository

func GetGroupById(w http.ResponseWriter, r *http.Request) {
	// Get group id from parameters
	params := mux.Vars(r)
	groupIDStr := params["groupId"]
	groupID, err := strconv.ParseUint(groupIDStr, 0, 32)
	if err != nil {
		ResponseError(w, fmt.Sprintf("Cannot parse group ID '%s': %s", groupIDStr, err), http.StatusBadRequest)
		return
	}

	// Get group from database
	group, err := GroupRepository.GetById(uint(groupID))
	if err != nil {
		ResponseError(w, fmt.Sprintf("Failed to get group by ID '%d'", groupID), http.StatusInternalServerError)
		return
	}

	// Write to response
	groupDTO := dto.GroupToGroupDTO(*group)
	ResponseJSON(w, groupDTO)
}

func GetGroupsByUser(w http.ResponseWriter, r *http.Request) {
	// Get user id from parameters
	params := mux.Vars(r)
	userIDStr := params["userId"]
	userID, err := strconv.ParseUint(userIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse user id '%s': %s", userIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Get groups from database
	groups, err := GroupRepository.GetByUser(uint(userID))
	if err != nil {
		errMsg := fmt.Sprintf("Error getting groups by user: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}

	// Write to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var groupDTOs []dto.GroupDTO
	for _, group := range groups {
		groupDTOs = append(groupDTOs, dto.GroupToGroupDTO(*group))
	}
	err = json.NewEncoder(w).Encode(groupDTOs)
	if err != nil {
		errMsg := fmt.Sprintf("Error encoding to json: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	// Parse group data from request body
	groupCreationRequest := &request.GroupCreationRequest{}
	err := json.NewDecoder(r.Body).Decode(groupCreationRequest)
	if err != nil {
		ResponseError(w, fmt.Sprintf("Cannot parse request body: %s", err), http.StatusBadRequest)
		return
	}
	group := &model.Group{
		Name: groupCreationRequest.Name,
	}
	// Get creator id from header
	userIDStr := r.Header.Get("userID")
	userID, _ := strconv.ParseUint(userIDStr, 0, 32)

	// Create group in database
	err = GroupRepository.Create(group, uint(userID))
	if err != nil {
		ResponseError(w, "Error creating group", http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Write to response
	groupDTO := dto.GroupToGroupDTO(*group)
	groupDTO.Members = nil
	groupDTO.Expenses = nil
	ResponseJSON(w, groupDTO)
}

func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	// Get group ID from parameters
	params := mux.Vars(r)
	groupIDStr := params["groupId"]
	groupID, err := strconv.ParseUint(groupIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse group ID '%s': %s", groupIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Parse group data from request body
	groupDTO := &dto.GroupDTO{}
	err = json.NewDecoder(r.Body).Decode(groupDTO)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse request body: %s", err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}
	groupDTO.ID = uint(groupID)
	group, err := groupDTO.MapToDomain()
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse model: %s", err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Update group to database
	err = GroupRepository.Update(&group)
	if err != nil {
		errMsg := fmt.Sprintf("Error while updating group: %s", err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Write updated data to response
	w.WriteHeader(http.StatusOK)
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	// Get user id from parameters
	params := mux.Vars(r)
	groupIDStr := params["groupId"]
	groupID, err := strconv.ParseUint(groupIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse group ID '%s': %s", groupIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Delete user from database and write response
	err = GroupRepository.Delete(uint(groupID))
	if err != nil {
		errMsg := fmt.Sprintf("Error deleting group with ID '%d': %s", groupID, err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}
	w.WriteHeader(http.StatusOK)
}
