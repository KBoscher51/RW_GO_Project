package handlers

import (
	"log"
	"net/http"
	"rw_local_go/db"
	"rw_local_go/models"
	"rw_local_go/repositories"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateVictime(c *gin.Context) {
	log.Println("=== Début CreateVictime ===")

	var victime models.Victime
	if err := c.ShouldBindJSON(&victime); err != nil {
		log.Printf("Erreur binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	victime.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	log.Printf("Victime à créer: %+v", victime)

	victimeRepository := repositories.VictimeRepository{DB: db.DB}

	if err := victimeRepository.Create(&victime); err != nil {
		log.Printf("Erreur création victime: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Victime créée avec succès")
	c.JSON(http.StatusCreated, gin.H{
		"message": "Victime créée",
		"victime": victime,
	})
}

func GetFiles(c *gin.Context) {
	log.Println("=== Début GetFiles ===")

	// Récupérer les paramètres d'authentification
	name := c.Query("name")
	password := c.Query("password")

	if name == "" || password == "" {
		log.Println("Paramètres d'authentification manquants")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nom et mot de passe requis"})
		return
	}

	log.Printf("Authentification pour: %s", name)

	// Authentifier la victime
	victimeRepository := repositories.VictimeRepository{DB: db.DB}
	victime, err := victimeRepository.Authenticate(name, password)
	if err != nil {
		log.Printf("Erreur authentification: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentification échouée"})
		return
	}

	// Récupérer les fichiers de cette victime
	fileRepository := repositories.FileRepository{DB: db.DB}
	files, err := fileRepository.GetFilesByVictimeID(victime.ID)
	if err != nil {
		log.Printf("Erreur récupération fichiers: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Fichiers trouvés: %d", len(files))
	c.JSON(http.StatusOK, gin.H{
		"victime": victime,
		"files":   files,
	})
}

func DownloadFile(c *gin.Context) {
	log.Println("=== Début DownloadFile ===")

	fileID := c.Param("id")
	name := c.Query("name")
	password := c.Query("password")

	if name == "" || password == "" {
		log.Println("Paramètres d'authentification manquants")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nom et mot de passe requis"})
		return
	}

	// Authentifier la victime
	victimeRepository := repositories.VictimeRepository{DB: db.DB}
	victime, err := victimeRepository.Authenticate(name, password)
	if err != nil {
		log.Printf("Erreur authentification: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentification échouée"})
		return
	}

	// Récupérer le fichier
	fileRepository := repositories.FileRepository{DB: db.DB}
	file, err := fileRepository.GetFileByID(fileID, victime.ID)
	if err != nil {
		log.Printf("Erreur récupération fichier: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Fichier non trouvé"})
		return
	}

	// Envoyer le fichier
	log.Printf("Envoi du fichier: %s", file.FilePath)
	c.File(file.FilePath)
}
