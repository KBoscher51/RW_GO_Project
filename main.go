package main

import (
	"rw_local_go/db"
	"rw_local_go/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDB()

	r := gin.Default()

	// Endpoints pour les fichiers
	r.POST("/upload", handlers.AddFile)
	r.GET("/files", handlers.GetFiles)
	r.GET("/download/:id", handlers.DownloadFile)

	// Endpoints pour les victimes
	r.POST("/victime", handlers.CreateVictime)

	r.Run(":8080")
}
