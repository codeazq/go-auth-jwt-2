package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type ApiError struct {
	Error string `json:"error"`
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/accounts", makeHTTPHandleFunc(s.handleGetAccounts)).Methods("GET")
	router.HandleFunc("/accounts", makeHTTPHandleFunc(s.handleCreateAccount)).Methods("POST")

	router.HandleFunc("/accounts/{id}", makeHTTPHandleFunc(s.handleDeleteAccount)).Methods("DELETE")
	router.HandleFunc("/accounts/{id}", makeHTTPHandleFunc(s.handleGetAccountById)).Methods("GET")
	router.HandleFunc("/accounts/{id}", makeHTTPHandleFunc(s.handleUpdateAccount)).Methods("PUT")

	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccounts(w, r)
	}

	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("method not supported %s", r.Method)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)

}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getParamterFromRequest(r, "id")

	if err != nil {
		return err
	}

	account, err := s.store.GetAccountById(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getParamterFromRequest(r, "id")

	if err != nil {
		return err
	}

	updateAccountRequest := new(UpdateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(updateAccountRequest); err != nil {
		return err
	}

	account, err := s.store.GetAccountById(id)
	if err != nil {
		return err
	}

	account.FirstName = updateAccountRequest.FirstName
	account.LastName = updateAccountRequest.LastName

	if err := s.store.UpdateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getParamterFromRequest(r, "id")

	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getParamterFromRequest(r *http.Request, parameter string) (int, error) {
	paramAsString := mux.Vars(r)[parameter]
	formatedParam, err := strconv.Atoi(paramAsString)

	if err != nil {
		return formatedParam, fmt.Errorf("error while parsing %s: %s", parameter, paramAsString)
	}

	return formatedParam, nil
}
