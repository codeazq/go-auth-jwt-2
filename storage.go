package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	GetAccountById(int) (*Account, error)
	UpdateAccount(*Account) error
	DeleteAccount(int) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "postgresql://postgres:postgres@localhost:5432/go_fintech_bank?sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(55),
		last_name VARCHAR(55),
		number SERIAL,
		balance DECIMAL(10, 2),
		created_at TIMESTAMP
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	sqlStatement := `INSERT INTO accounts 
	(first_name, last_name, number, balance, created_at)
	VALUES 
	($1, $2, $3, $4, $5)`

	log.Println("before the query")
	resp, err := s.db.Query(
		sqlStatement, acc.FirstName, acc.LastName,
		acc.Number, acc.Balance, acc.CreatedAt)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)

	return nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	return nil, nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}
