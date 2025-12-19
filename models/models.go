package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FullName  string    `json:"full_name"`
	HandID    string    `json:"hand_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Wallet struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id"`
	Currency    string    `json:"currency"`
	Balance     float64   `json:"balance"`
	LastUpdated time.Time `json:"last_updated"`
}

type Transaction struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	WalletID  uint      `json:"wallet_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`

	// --- CHAMPS BLOCKCHAIN ---
	PreviousHash string    `json:"previous_hash"`
	Hash         string    `json:"hash"`
	HandToken    string    `json:"hand_token"`

	CreatedAt time.Time `json:"created_at"`
}