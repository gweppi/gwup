package main


import (
	"fmt"
	"net/http"
	"io"
	"os"
	"mime"
	"time"

	"github.com/gweppi/gwup/internal/shared"
	"github.com/gweppi/gwup/cmd/server/utils"
	"github.com/gin-gonic/gin"
)

const DATA_DIR = "./data/"
const MiB = 1 << 20

func main() {
	router := gin.Default()

	router.GET("/health", handleHealth)
	router.POST("/upload", handleUpload)
	router.GET("/download", handleDownload)

	router.Run()
}

func handleHealth(ctx *gin.Context) {
	serverInfo := shared.ServerInfo{
		Status: "ok",
		Version: "1.0.0",
		RequiresAuth: false,
	}

	ctx.IndentedJSON(http.StatusOK, serverInfo)
}

func handleUpload(ctx *gin.Context) {
	if ctx.Request.Body == nil {
		ctx.Status(400)
		return
	}

	maxReader := http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 100 * MiB)
	_, params, err := mime.ParseMediaType(ctx.Request.Header.Get("Content-Disposition"))
	if err != nil {
		ctx.Status(500)
		fmt.Println(err)
		return
	}

	fileName := params["filename"]
	if fileName == "" {
		ctx.Status(500)
		fmt.Println(err)
		return
	}

	file, err := os.OpenFile(DATA_DIR + fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		ctx.Status(500)
		fmt.Println(err)
		return
	}

	if _, err := io.Copy(file, maxReader); err != nil {
		ctx.Status(500)
		fmt.Println(err)
		return
	}

	ctx.Status(201)
}

func handleDownload(ctx *gin.Context) {
	fileId := ctx.Request.Header.Get("X-File-Id")

	entries, err := os.ReadDir(DATA_DIR)
	if err != nil {
		// Error reading dir
		ctx.Status(500)
		return
	}

	fileName := ""
	var time time.Time
	for _, entry := range entries {
		if fileId != "" {
			if utils.FileNameMatches(entry.Name(), fileId) {
				fileName = entry.Name()
				break
			}
		} else {
			fileInfo, err := entry.Info()
			if err != nil {
				ctx.Status(500)
				return
			}

			if time.IsZero() || fileInfo.ModTime().Compare(time) == 1 {
				fileName = entry.Name()
				time = fileInfo.ModTime()
			} 
		}
	}
	
	if fileName == "" {
		// File was not found in directory
		ctx.Status(400)
		return
	}

	file, err := os.OpenFile(DATA_DIR + fileName, os.O_RDONLY, 0644)
	if err != nil {
		// Error opening file
		ctx.Status(500)
		return
	}
	defer file.Close()

	fileStats, err := file.Stat()
	if err != nil {
		// Something went wrong with parsing into multipart
		ctx.Status(500)
		return
	}

	extraHeaders := map[string]string {
		"Content-Disposition": "attachment; filename=\"" + fileName + "\"",
	}
	ctx.DataFromReader(http.StatusOK, fileStats.Size(), "application/octet-stream", file, extraHeaders)
}

