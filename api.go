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
	Error string
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

	router.HandleFunc("/accounts", makeHTTPHandleFunc(s.handleAccount))

	router.HandleFunc("/accounts/{id}", makeHTTPHandleFunc(s.handleGetAccountById))

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
	accountId := mux.Vars(r)["id"]
	ID, parsingErr := strconv.Atoi(accountId)

	if parsingErr != nil {
		fmt.Println("error while parsing")
		return parsingErr
	}

	account, err := s.store.GetAccountById(ID)

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

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
