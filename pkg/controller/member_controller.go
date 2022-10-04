package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"money_share/pkg/dto"
	"money_share/pkg/repository"
	"net/http"
	"strconv"
)

var MemberRepository repository.MemberRepository

func GetMemberByID(w http.ResponseWriter, r *http.Request) {
	// Get user id and group id form query params
	queries := r.URL.Query()
	userIDStr := queries.Get("userId")
	groupIDStr := queries.Get("groupId")
	if len(groupIDStr) == 0 || len(groupIDStr) == 0 {
		errMsg := fmt.Sprintf("Not enough parameters provided, required 'userId' and 'groupId'")
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse user ID '%s': %s", userIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}
	groupID, err := strconv.ParseUint(groupIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse group ID '%s': %s", groupIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Get member from database
	member, err := MemberRepository.GetByID(uint(userID), uint(groupID))
	if err != nil {
		errMsg := fmt.Sprintf("Error getting member by user '%d', group '%d': %s", userID, groupID, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Write to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	memberDTO := dto.MemberToMemberDTO(*member)
	err = json.NewEncoder(w).Encode(memberDTO)
	if err != nil {
		errMsg := fmt.Sprintf("Error encoding to json: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}
}

func GetMembersOfGroup(w http.ResponseWriter, r *http.Request) {
	// Get group id from parameters
	params := mux.Vars(r)
	groupIDStr := params["groupId"]
	groupID, err := strconv.ParseUint(groupIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse group id '%s': %s", groupIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Get members from database
	members, err := MemberRepository.GetByGroup(uint(groupID))
	if err != nil {
		errMsg := fmt.Sprintf("Error getting members by group: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}

	// Write to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var memberDTOs []dto.MemberDTO
	for _, member := range members {
		memberDTOs = append(memberDTOs, dto.MemberToMemberDTO(*member))
	}
	err = json.NewEncoder(w).Encode(memberDTOs)
	if err != nil {
		errMsg := fmt.Sprintf("Error encoding to json: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}
}

func AddMemberToGroup(w http.ResponseWriter, r *http.Request) {
	// Get user id and group id form query params
	queries := r.URL.Query()
	userIDStr := queries.Get("userId")
	groupIDStr := queries.Get("groupId")
	if len(userIDStr) == 0 || len(groupIDStr) == 0 {
		errMsg := fmt.Sprintf("Not enough parameters provided, required 'userId' and 'groupId'")
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse user ID '%s': %s", userIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}
	groupID, err := strconv.ParseUint(groupIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse group ID '%s': %s", groupIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Add member to group in database
	err = MemberRepository.AddMemberToGroup(uint(userID), uint(groupID))
	if err != nil {
		errMsg := fmt.Sprintf("Error adding member to group: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}

	// Write to response
	w.WriteHeader(http.StatusOK)
}

func RemoveMemberFromGroup(w http.ResponseWriter, r *http.Request) {
	// Get user id and group id form query params
	queries := r.URL.Query()
	userIDStr := queries.Get("userId")
	groupIDStr := queries.Get("groupId")
	if len(groupIDStr) == 0 || len(groupIDStr) == 0 {
		errMsg := fmt.Sprintf("Not enough parameters provided, required 'userId' and 'groupId'")
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse user ID '%s': %s", userIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}
	groupID, err := strconv.ParseUint(groupIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse group ID '%s': %s", groupIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Remove member from group in database
	err = MemberRepository.RemoveMemberFromGroup(uint(userID), uint(groupID))
	if err != nil {
		errMsg := fmt.Sprintf("Error removing member from group: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}

	// Write to response
	w.WriteHeader(http.StatusOK)
}
