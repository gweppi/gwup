package main


import (
	"fmt"
	"net/http"

	"github.com/gweppi/gwup/internal/shared"
	"github.com/gin-gonic/gin"
)


func main() {
	fmt.Println("Hello, World!")

	router := gin.Default()

	router.GET("/health", handleHealth)

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
