package repositories

import (
	"log"
	"rw_local_go/models"

	"github.com/jmoiron/sqlx"
)

type VictimeRepository struct {
	DB *sqlx.DB
}

func (v *VictimeRepository) Authenticate(name, password string) (*models.Victime, error) {
	log.Printf("Authentification pour: %s", name)

	var victime models.Victime
	query := `SELECT id, name, password, created_at FROM victimes WHERE name = ? AND password = ?`

	err := v.DB.Get(&victime, query, name, password)
	if err != nil {
		log.Printf("Erreur authentification: %v", err)
		return nil, err
	}

	log.Printf("Authentification réussie pour: %s", name)
	return &victime, nil
}

func (v *VictimeRepository) Create(victime *models.Victime) error {
	log.Printf("Création victime: %s", victime.Name)

	query := `INSERT INTO victimes (name, password, created_at) VALUES (?, ?, ?)`

	result, err := v.DB.Exec(query, victime.Name, victime.Password, victime.CreatedAt)
	if err != nil {
		log.Printf("Erreur création victime: %v", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Erreur récupération ID: %v", err)
		return err
	}

	victime.ID = int(id)
	log.Printf("Victime créée avec ID: %d", victime.ID)
	return nil
}
