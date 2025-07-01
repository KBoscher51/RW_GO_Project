package repositories

import (
	"log"
	"rw_local_go/models"

	"github.com/jmoiron/sqlx"
)

type FileRepository struct {
	DB *sqlx.DB
}

func (file *FileRepository) Create(fileModel *models.File) error {
	log.Println("=== Début Create ===")
	log.Printf("DB pointer dans Create: %v", file.DB)
	log.Printf("FileModel: %+v", fileModel)

	query := `INSERT INTO files (filename, filepath, filepathvictime, victime_id, filesize, mimetype, created_at) 
			  VALUES (:FileName, :FilePath, :FilePathVictime, :VictimeID, :FileSize, :MimeType, :CreatedAt)`
	log.Printf("Query: %s", query)

	_, err := file.DB.NamedExec(query, map[string]interface{}{
		"FileName":        fileModel.FileName,
		"FilePath":        fileModel.FilePath,
		"FilePathVictime": fileModel.FilePathVictime,
		"VictimeID":       fileModel.VictimeID,
		"FileSize":        fileModel.FileSize,
		"MimeType":        fileModel.MimeType,
		"CreatedAt":       fileModel.CreatedAt,
	})
	if err != nil {
		log.Printf("Erreur NamedExec: %v", err)
		return err
	}

	log.Println("Insertion réussie")
	return nil
}

func (file *FileRepository) GetFilesByVictimeID(victimeID int) ([]models.File, error) {
	log.Printf("Récupération fichiers pour victime ID: %d", victimeID)

	var files []models.File
	query := `SELECT idfiles, filename, filepath, filepathvictime, victime_id, filesize, mimetype, created_at 
			  FROM files WHERE victime_id = ? ORDER BY created_at DESC`

	err := file.DB.Select(&files, query, victimeID)
	if err != nil {
		log.Printf("Erreur récupération fichiers: %v", err)
		return nil, err
	}

	log.Printf("Fichiers trouvés: %d", len(files))
	return files, nil
}

func (file *FileRepository) GetFileByID(fileID string, victimeID int) (*models.File, error) {
	log.Printf("Récupération fichier ID: %s pour victime ID: %d", fileID, victimeID)

	var fileModel models.File
	query := `SELECT idfiles, filename, filepath, filepathvictime, victime_id, filesize, mimetype, created_at 
			  FROM files WHERE idfiles = ? AND victime_id = ?`

	err := file.DB.Get(&fileModel, query, fileID, victimeID)
	if err != nil {
		log.Printf("Erreur récupération fichier: %v", err)
		return nil, err
	}

	log.Printf("Fichier trouvé: %s", fileModel.FileName)
	return &fileModel, nil
}
