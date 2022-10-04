package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"money_share/pkg/dto"
	"money_share/pkg/http/request"
	"money_share/pkg/repository"
	"net/http"
	"strconv"
)

var ExpenseRepository repository.ExpenseRepository

func GetExpenseByID(w http.ResponseWriter, r *http.Request) {
	// Get expense id from parameters
	params := mux.Vars(r)
	expenseIDStr := params["expenseId"]
	expenseID, err := strconv.ParseUint(expenseIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse expense id '%s': %s", expenseIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Get expense from database
	expense, err := ExpenseRepository.GetById(uint(expenseID))
	if err != nil {
		errMsg := fmt.Sprintf("Error getting expense by id '%d': %s", expenseID, err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}

	// Write to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	expenseDTO := dto.ExpenseToExpenseDTO(*expense)
	err = json.NewEncoder(w).Encode(expenseDTO)
	if err != nil {
		errMsg := fmt.Sprintf("Error encoding to json: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}
}

func GetExpensesByGroup(w http.ResponseWriter, r *http.Request) {
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

	// Get expenses from database
	expenses, err := ExpenseRepository.GetByGroup(uint(groupID))
	if err != nil {
		errMsg := fmt.Sprintf("Error getting expense by group id '%d': %s", groupID, err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}

	// Write to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var expenseDTOs []dto.ExpenseDTO
	for _, expense := range expenses {
		expenseDTOs = append(expenseDTOs, dto.ExpenseToExpenseDTO(*expense))
	}
	err = json.NewEncoder(w).Encode(expenseDTOs)
	if err != nil {
		errMsg := fmt.Sprintf("Error encoding to json: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}
}

func GetExpensesByMember(w http.ResponseWriter, r *http.Request) {
	// Get member id and group id from queries
	queries := r.URL.Query()
	memberIDStr := queries.Get("memberId")
	groupIDStr := queries.Get("groupId")
	if len(memberIDStr) == 0 || len(groupIDStr) == 0 {
		errMsg := fmt.Sprintf("Not enough parameters provided, required 'memberId' and 'groupId'")
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}
	memberID, err := strconv.ParseUint(memberIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse member id '%s': %s", memberIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}
	groupID, err := strconv.ParseUint(groupIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse group id '%s': %s", groupIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Get expenses from database
	expenses, err := ExpenseRepository.GetByMember(uint(memberID), uint(groupID))
	if err != nil {
		errMsg := fmt.Sprintf("Error getting expenses: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}

	// Write to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var expenseDTOs []dto.ExpenseDTO
	for _, expense := range expenses {
		expenseDTOs = append(expenseDTOs, dto.ExpenseToExpenseDTO(*expense))
	}
	err = json.NewEncoder(w).Encode(expenseDTOs)
	if err != nil {
		errMsg := fmt.Sprintf("Error encoding to json: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}
}

func CreateExpense(w http.ResponseWriter, r *http.Request) {
	// Parse expense creation request from request body
	expenseCreationRequest := &request.ExpenseCreationRequest{}
	err := json.NewDecoder(r.Body).Decode(expenseCreationRequest)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse request body: %s", err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}
	// TODO: Validate fields
	expense := expenseCreationRequest.Expense.MapToDomain()
	expense.MemberID = expenseCreationRequest.MemberID
	expense.GroupID = expenseCreationRequest.GroupID

	// Create expense in database
	err = ExpenseRepository.Create(&expense)
	if err != nil {
		errMsg := fmt.Sprintf("Error creating expense: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}

	// Write to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	expenseDTO := dto.ExpenseToExpenseDTO(expense)
	err = json.NewEncoder(w).Encode(expenseDTO)
	if err != nil {
		errMsg := fmt.Sprintf("Error encoding to json: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}
}

func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	// Get expense id from parameters
	params := mux.Vars(r)
	expenseIDStr := params["expenseId"]
	expenseID, err := strconv.ParseUint(expenseIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse expense id '%s': %s", expenseIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Parse expense data from request body
	expenseDTO := &dto.ExpenseDTO{}
	err = json.NewDecoder(r.Body).Decode(expenseDTO)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse request body: %s", err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}
	expense := expenseDTO.MapToDomain()
	expense.ID = uint(expenseID)

	// Update expense in database
	err = ExpenseRepository.Update(&expense)
	if err != nil {
		errMsg := fmt.Sprintf("Error updating expense: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}

	// Write to response
	w.WriteHeader(http.StatusOK)
}

func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	// Get expense id from parameters
	params := mux.Vars(r)
	expenseIDStr := params["expenseId"]
	expenseID, err := strconv.ParseUint(expenseIDStr, 0, 32)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse expense id '%s': %s", expenseIDStr, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		fmt.Println(errMsg)
		return
	}

	// Delete expense from database
	err = ExpenseRepository.Delete(uint(expenseID))
	if err != nil {
		errMsg := fmt.Sprintf("Error deleting expense: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		fmt.Println(errMsg)
		return
	}

	// Write to response
	w.WriteHeader(http.StatusOK)
}
