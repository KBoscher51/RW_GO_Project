package models

type File struct {
	ID              string `json:"id" db:"idfiles"`
	FileName        string `json:"FileName" db:"filename"`
	FilePath        string `json:"FilePath" db:"filepath"`               // Chemin sur le disque
	FilePathVictime string `json:"FilePathVictime" db:"filepathvictime"` // Chemin local transmis
	VictimeID       int    `json:"VictimeID" db:"victime_id"`            // ID de la victime
	FileSize        int64  `json:"FileSize" db:"filesize"`               // Taille en bytes
	MimeType        string `json:"MimeType" db:"mimetype"`               // Type MIME
	CreatedAt       string `json:"CreatedAt" db:"created_at"`
}
