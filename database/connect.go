package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// ---------------------------------------------------------
	// COLLE TA PHRASE NEON JUSTE EN DESSOUS ENTRE LES GUILLEMETS
	// ---------------------------------------------------------
	dsn := "postgresql://neondb_owner:npg_PH0sn4WLAyBk@ep-withered-term-agwunzo0-pooler.c-2.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("❌ Erreur critique : Impossible de se connecter au Cloud !", err)
	}

	fmt.Println("✅ Succès : Connecté à la Banque Cloud (Neon/AWS) ☁️")
}