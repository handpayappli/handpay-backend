package main

import (
	"fmt"
	"handpay/database"
	"handpay/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Structure pour recevoir les données du paiement
type PaymentRequest struct {
	UserID    uint    `json:"user_id"`
	Amount    float64 `json:"amount"`
	HandToken string  `json:"hand_token"`
}

func main() {
	fmt.Println("🚀 Démarrage du système HandPay v1.0 (Cloud Edition)...")

	// 1. Connexion BDD & Migration
	database.Connect()
	database.DB.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Transaction{})

	// 2. App Web
	app := fiber.New()

	// IMPORTANT : Autoriser les connexions depuis le téléphone (CORS)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Tout le monde peut se connecter (pour le test)
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("HandPay System is Online ☁️")
	})

	// --- ROUTE 1 : SOLDE ---
	app.Get("/balance/:user_id", func(c *fiber.Ctx) error {
		userId := c.Params("user_id")
		var wallet models.Wallet
		if err := database.DB.Where("user_id = ?", userId).First(&wallet).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Wallet introuvable"})
		}
		return c.JSON(fiber.Map{"user_id": wallet.UserID, "balance": wallet.Balance, "currency": wallet.Currency})
	})

	// --- ROUTE 2 : PAIEMENT SECURISE ---
	app.Post("/pay", func(c *fiber.Ctx) error {
		var req PaymentRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Données invalides"})
		}

		fmt.Printf("\n🕵️  DEBUG REÇU -> UserID: %d | Montant: %.2f | Token: '%s'\n", req.UserID, req.Amount, req.HandToken)

		tx := database.DB.Begin()

		var wallet models.Wallet
		if err := tx.Where("user_id = ?", req.UserID).First(&wallet).Error; err != nil {
			tx.Rollback()
			return c.Status(404).JSON(fiber.Map{"error": "Wallet introuvable"})
		}
		if wallet.Balance < req.Amount {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Solde insuffisant"})
		}

		wallet.Balance -= req.Amount
		if err := tx.Save(&wallet).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Erreur débit"})
		}

		// Blockchain logic
		var lastTx models.Transaction
		database.DB.Last(&lastTx)
		prevHash := lastTx.Hash
		if prevHash == "" {
			prevHash = "GENESIS"
		}

		newHash := models.CalculateHash(wallet.ID, -req.Amount, prevHash, req.HandToken)

		transaction := models.Transaction{
			WalletID:     wallet.ID,
			Amount:       -req.Amount,
			Type:         "PAYMENT",
			Status:       "SUCCESS",
			PreviousHash: prevHash,
			Hash:         newHash,
			HandToken:    req.HandToken,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Erreur Blockchain"})
		}

		tx.Commit()

		return c.JSON(fiber.Map{
			"message":     "Paiement Validé & Tokenisé 🔒",
			"tx_hash":     newHash,
			"new_balance": wallet.Balance,
		})
	})

	// --- ROUTE 3 : INSCRIPTION (Nouveau) ---
	app.Post("/register", func(c *fiber.Ctx) error {
		type RegisterRequest struct {
			Name string `json:"name"`
		}
		var req RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Nom requis"})
		}

		// Création Client
		user := models.User{FullName: req.Name, HandID: "temp-" + req.Name}
		database.DB.Create(&user)

		// Création Wallet avec 500€
		wallet := models.Wallet{UserID: user.ID, Balance: 500.00, Currency: "EUR"}
		database.DB.Create(&wallet)

		return c.JSON(fiber.Map{
			"message": "Client Cloud créé !",
			"user_id": user.ID,
			"balance": wallet.Balance,
		})
	})

		fmt.Println("👂 En écoute sur toutes les interfaces (0.0.0.0:3000)")
	app.Listen("0.0.0.0:3000")
}