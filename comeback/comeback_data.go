package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	API_BASE_URL = "http://localhost:8080"
	VICTIM_NAME  = "John"
	PASSWORD     = "TEST123"
)

type Victime struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type File struct {
	ID              string `json:"id"`
	FileName        string `json:"FileName"`
	FilePath        string `json:"FilePath"`
	FilePathVictime string `json:"FilePathVictime"`
	VictimeID       int    `json:"VictimeID"`
	FileSize        int64  `json:"FileSize"`
	MimeType        string `json:"MimeType"`
	CreatedAt       string `json:"CreatedAt"`
}

type FilesResponse struct {
	Victime Victime `json:"victime"`
	Files   []File  `json:"files"`
}

type DataRestorer struct {
	client *http.Client
}

func NewDataRestorer() *DataRestorer {
	return &DataRestorer{
		client: &http.Client{},
	}
}

func (dr *DataRestorer) getFiles() (*FilesResponse, error) {
	fmt.Println("📋 Récupération de la liste des fichiers...")

	// Construire l'URL avec les paramètres d'authentification
	url := fmt.Sprintf("%s/files?name=%s&password=%s", API_BASE_URL, VICTIM_NAME, PASSWORD)

	// Créer la requête
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erreur création requête: %v", err)
	}

	// Envoyer la requête
	resp, err := dr.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erreur envoi requête: %v", err)
	}
	defer resp.Body.Close()

	// Lire la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture réponse: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur récupération fichiers: %s - %s", resp.Status, string(body))
	}

	// Parser la réponse JSON
	var filesResp FilesResponse
	if err := json.Unmarshal(body, &filesResp); err != nil {
		return nil, fmt.Errorf("erreur parsing JSON: %v", err)
	}

	fmt.Printf("%d fichiers trouvés pour %s\n", len(filesResp.Files), filesResp.Victime.Name)
	return &filesResp, nil
}

func (dr *DataRestorer) downloadFile(fileID, originalPath string) error {
	fmt.Printf("Téléchargement: %s\n", filepath.Base(originalPath))

	// Construire l'URL de téléchargement
	url := fmt.Sprintf("%s/download/%s?name=%s&password=%s", API_BASE_URL, fileID, VICTIM_NAME, PASSWORD)

	// Créer la requête
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("erreur création requête: %v", err)
	}

	// Envoyer la requête
	resp, err := dr.client.Do(req)
	if err != nil {
		return fmt.Errorf("erreur envoi requête: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erreur téléchargement: %s", resp.Status)
	}

	// Créer le dossier parent si nécessaire
	dir := filepath.Dir(originalPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erreur création dossier %s: %v", dir, err)
	}

	// Créer le fichier de destination
	destFile, err := os.Create(originalPath)
	if err != nil {
		return fmt.Errorf("erreur création fichier %s: %v", originalPath, err)
	}
	defer destFile.Close()

	// Copier le contenu
	_, err = io.Copy(destFile, resp.Body)
	if err != nil {
		return fmt.Errorf("erreur copie fichier: %v", err)
	}

	fmt.Printf("Fichier restauré: %s\n", originalPath)
	return nil
}

func (dr *DataRestorer) restoreData() error {
	fmt.Println("Début de la restauration des données...")
	fmt.Printf("Victime: %s\n", VICTIM_NAME)
	fmt.Printf("API: %s\n", API_BASE_URL)

	// Récupérer la liste des fichiers
	filesResp, err := dr.getFiles()
	if err != nil {
		return fmt.Errorf("erreur récupération fichiers: %v", err)
	}

	if len(filesResp.Files) == 0 {
		fmt.Println("Aucun fichier à restaurer")
		return nil
	}

	// Statistiques
	successCount := 0
	errorCount := 0

	fmt.Printf("\nRestauration de %d fichiers...\n", len(filesResp.Files))

	// Restaurer chaque fichier
	for _, file := range filesResp.Files {
		// Utiliser le chemin original de la victime
		originalPath := file.FilePathVictime

		// Convertir les chemins Windows si nécessaire
		originalPath = strings.ReplaceAll(originalPath, "/", "\\")

		// Télécharger et restaurer le fichier
		if err := dr.downloadFile(file.ID, originalPath); err != nil {
			fmt.Printf("Erreur restauration %s: %v\n", filepath.Base(originalPath), err)
			errorCount++
		} else {
			successCount++
		}
	}

	// Résumé
	fmt.Printf("\nRestauration terminée!\n")
	fmt.Printf("Fichiers restaurés avec succès: %d\n", successCount)
	if errorCount > 0 {
		fmt.Printf("Fichiers en erreur: %d\n", errorCount)
	}

	return nil
}

func main() {
	fmt.Println("COMEBACK DATA - Restauration des fichiers")
	fmt.Println("=============================================")

	restorer := NewDataRestorer()

	if err := restorer.restoreData(); err != nil {
		log.Fatalf("Erreur: %v", err)
	}

	fmt.Println("\nRestauration terminée avec succès!")
}
