package db

import (
	"fmt"
	"log"
	"rw_local_go/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func ConnectDB() {
	log.Println("=== Début ConnectDB ===")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Erreur chargement config: %v", err)
		log.Fatal(err)
	}

	log.Printf("Config chargée: %+v", cfg)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.DB_USER,
		cfg.DB_PASSWORD,
		cfg.DB_HOST,
		cfg.DB_PORT,
		cfg.DB_NAME,
	)

	log.Printf("DSN: %s", dsn)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Printf("Erreur connexion DB: %v", err)
		log.Fatal(err)
	}

	log.Println("Connexion DB réussie")
	DB = db
}
