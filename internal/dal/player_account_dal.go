package dal

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"mud/internal/models"
)

// PlayerAccountDALInterface defines the interface for player account data operations.
type PlayerAccountDALInterface interface {
	CreateAccount(username, password, email string) (*models.PlayerAccount, error)
	GetAccountByUsername(username string) (*models.PlayerAccount, error)
	Authenticate(username, password string) (*models.PlayerAccount, error)
	UpdateLastLogin(accountID string) error
}

// PlayerAccountDAL implements the PlayerAccountDALInterface.
type PlayerAccountDAL struct {
	DB *sql.DB
}

// NewPlayerAccountDAL creates a new PlayerAccountDAL.
func NewPlayerAccountDAL(db *sql.DB) *PlayerAccountDAL {
	return &PlayerAccountDAL{DB: db}
}

// CreateAccount creates a new player account.
func (dal *PlayerAccountDAL) CreateAccount(username, password, email string) (*models.PlayerAccount, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	account := &models.PlayerAccount{
		ID:             uuid.New().String(),
		Username:       username,
		HashedPassword: string(hashedPassword),
		Email:          email,
		CreatedAt:      time.Now(),
	}

	query := `INSERT INTO player_accounts (id, username, hashed_password, email, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err = dal.DB.Exec(query, account.ID, account.Username, account.HashedPassword, account.Email, account.CreatedAt)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetAccountByUsername retrieves a player account by username.
func (dal *PlayerAccountDAL) GetAccountByUsername(username string) (*models.PlayerAccount, error) {
	query := `SELECT id, username, hashed_password, email, created_at, last_login_at FROM player_accounts WHERE username = ?`
	row := dal.DB.QueryRow(query, username)

	var account models.PlayerAccount
	var lastLogin sql.NullTime
	err := row.Scan(&account.ID, &account.Username, &account.HashedPassword, &account.Email, &account.CreatedAt, &lastLogin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil, nil if not found
		}
		return nil, err
	}
	if lastLogin.Valid {
		account.LastLoginAt = lastLogin.Time
	}

	return &account, nil
}

// Authenticate checks the username and password and returns the account if valid.
func (dal *PlayerAccountDAL) Authenticate(username, password string) (*models.PlayerAccount, error) {
	account, err := dal.GetAccountByUsername(username)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, nil // User not found
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.HashedPassword), []byte(password))
	if err != nil {
		return nil, nil // Password does not match
	}

	return account, nil
}

// UpdateLastLogin updates the last_login_at timestamp for a given account.
func (dal *PlayerAccountDAL) UpdateLastLogin(accountID string) error {
	query := `UPDATE player_accounts SET last_login_at = ? WHERE id = ?`
	_, err := dal.DB.Exec(query, time.Now(), accountID)
	if err != nil {
		log.Printf("Error updating last login for account %s: %v", accountID, err)
		return err
	}
	return nil
}
