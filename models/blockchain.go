package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func CalculateHash(walletID uint, amount float64, prevHash string, handToken string) string {
	record := fmt.Sprintf("%d%.2f%s%s%d", walletID, amount, prevHash, handToken, time.Now().UnixNano())
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}