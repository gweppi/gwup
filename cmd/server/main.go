package main


import (
	"fmt"
	"net/http"
	"io"
	"os"
	"mime"
	"mime/multipart"

	"github.com/gweppi/gwup/internal/shared"
	"github.com/gin-gonic/gin"
)


func main() {
	fmt.Println("Hello, World!")

	router := gin.Default()

	router.GET("/health", handleHealth)
	router.POST("/upload", handleUpload)

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
	// file, err := os.OpenFile("testfile.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	// if err != nil {
	// 	ctx.Status(500)
	// 	fmt.Println(err)
	// 	return
	// }
	// defer file.Close()
	// if _, err := io.Copy(file, ctx.Request.Body); err != nil {
	// 	ctx.Status(500)
	// 	fmt.Println("error copying file", err)
	// 	return
	// }

	_, params, err := mime.ParseMediaType(ctx.Request.Header.Get("Content-Type"))
	if err != nil {
		ctx.Status(500)
		fmt.Println(err)
		return
	}

	if boundary, ok := params["boundary"]; ok {
		multipartReader := multipart.NewReader(ctx.Request.Body, boundary)
		part, err := multipartReader.NextPart()
		if err == nil {
			defer part.Close()

			fmt.Println(part.Header)
			_, params, err := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
			if err != nil {
				fmt.Println("Something went wrong parsing Content-Disposition headers")
				return
			}
			filename, ok := params["filename"]
			if (!ok) {
				fmt.Println("Wrongly formatted request body")
				ctx.Status(400)
				return
			}

			file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				fmt.Println("Something went wrong creating a new file")
				return
			}
			defer file.Close()

			if _, err := io.Copy(file, part); err != nil {
				fmt.Println("Something went wrong writing to the file")
				return
			}

			ctx.Status(201)
			return
		} else {
			if err == io.EOF {
				fmt.Println("There are no more parts to read")
			} else {
				fmt.Println("error reading next part:", err)
			}
		}
	}
	ctx.Status(500)
}
