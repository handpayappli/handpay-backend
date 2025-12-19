package main

import (
	"fmt"
	"handpay/database"
	"handpay/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// DTOs (Structures de données pour les requêtes)
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MerchantPayRequest struct {
	MerchantID uint    `json:"merchant_id"`
	Amount     float64 `json:"amount"`
	HandToken  string  `json:"hand_token"` // Le token de la main du CLIENT
}

func main() {
	fmt.Println("🚀 Démarrage HandPay V2...")

	database.Connect()
	database.DB.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Transaction{})

	app := fiber.New()
	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))
	app.Static("/", "../frontend/app.html")

	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("HandPay V2 Online 🔒") })

	// --- 1. INSCRIPTION COMPLÈTE ---
	app.Post("/auth/register", func(c *fiber.Ctx) error {
		var req RegisterRequest
		if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(fiber.Map{"error": "Données invalides"}) }

		// Vérifier si email existe déjà
		var existing models.User
		if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
			return c.Status(400).JSON(fiber.Map{"error": "Email déjà utilisé"})
		}

		// Créer User (Dans la vraie vie, on hacherait le mot de passe avec bcrypt ici)
		user := models.User{FullName: req.Name, Email: req.Email, Password: req.Password, HandID: "pending"}
		database.DB.Create(&user)

		// Créer Wallet
		wallet := models.Wallet{UserID: user.ID, Balance: 0.00, Currency: "EUR"}
		database.DB.Create(&wallet)

		return c.JSON(fiber.Map{"message": "Compte créé", "user_id": user.ID})
	})

	// --- 2. CONNEXION ---
	app.Post("/auth/login", func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(fiber.Map{"error": "Données invalides"}) }

		var user models.User
		if err := database.DB.Where("email = ? AND password = ?", req.Email, req.Password).First(&user).Error; err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Email ou mot de passe incorrect"})
		}

		return c.JSON(fiber.Map{"message": "Connecté", "user_id": user.ID, "name": user.FullName})
	})

	// --- 3. LIER UNE MAIN (Enrôlement) ---
	app.Post("/user/link-hand", func(c *fiber.Ctx) error {
		type LinkHandReq struct {
			UserID    uint   `json:"user_id"`
			HandToken string `json:"hand_token"`
		}
		var req LinkHandReq
		if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(fiber.Map{"error": "Erreur"}) }

		// On met à jour le token de main de l'utilisateur
		// Note : Dans un vrai système, on stockerait ça dans Qdrant avec l'ID user.
		// Ici, on fait simple : on met à jour la BDD pour dire "Cet user a une main enregistrée".
		database.DB.Model(&models.User{}).Where("id = ?", req.UserID).Update("hand_id", "linked")
		
		return c.JSON(fiber.Map{"message": "Main liée au compte !"})
	})

	// --- 4. PAIEMENT COMMERÇANT (Scan de main Client) ---
	app.Post("/pay/merchant", func(c *fiber.Ctx) error {
		var req MerchantPayRequest
		if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(fiber.Map{"error": "Invalid"}) }

		// ÉTAPE A : IDENTIFIER LE CLIENT VIA SA MAIN
		// Ici, le backend devrait demander à Qdrant : "A qui est ce token ?"
		// Pour notre prototype V2 simplifié sans Qdrant Backend complexe :
		// On suppose que le Mobile envoie l'ID du client trouvé par l'IA, ou on simule.
		
		// ⚠️ Astuce pour le prototype : On va dire que le Scanner Python (Cloud AI) renvoie l'User ID.
		// Mais si c'est un NOUVEAU client, il faut qu'il ait déjà enregistré sa main.
		
		// Pour l'instant, utilisons la logique de paiement classique mais inversée (Client paie Commerçant)
		// On a besoin de savoir QUI est le client. C'est le rôle de l'IA.
		// On va supposer que l'appli mobile a déjà interrogé l'IA et récupéré l'ID du payeur.
		
		return c.Status(501).JSON(fiber.Map{"error": "Fonction en cours de dév - Besoin de lier IA et UserID"})
	})
	
	// Garde les anciennes routes pour compatibilité temporaire...
	app.Get("/balance/:user_id", func(c *fiber.Ctx) error {
		userId := c.Params("user_id")
		var wallet models.Wallet
		if err := database.DB.Where("user_id = ?", userId).First(&wallet).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Wallet introuvable"})
		}
		return c.JSON(fiber.Map{"balance": wallet.Balance})
	})

	fmt.Println("👂 V2 En écoute sur :3000")
	app.Listen(":3000")
}