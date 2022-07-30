package main

import (
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
	res := make(chan gin.H)
	requestCtx := ctx.Request.Context()

	go func() {
		time.Sleep(5 * time.Second)
		res <- gin.H{"msg": "Hello"}
		close(res)
	}()
	for {
		select {
		case dst := <-res:
			ctx.JSON(http.StatusOK, dst)
			return
		case <-requestCtx.Done():
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Error"})
			return
		}
	}

}
