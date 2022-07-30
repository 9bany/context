package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("", handle)
	router.Run()
}

func handle(ctx *gin.Context) {

	defer log.Println("i was sent payload to you")
	requestCtx := ctx.Request.Context()
	select {
	case <-time.After(5 * time.Second):
		ctx.JSON(http.StatusOK, gin.H{"msg": "Hello"})
	case <-requestCtx.Done():
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Error"})
	}
}
