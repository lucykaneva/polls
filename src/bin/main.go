package main

import (
	"log"
	handlers "polls/src/pckg/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/polls", handlers.GetAllPolls)
	router.GET("/polls/:id", handlers.GetAPollByID)
	router.POST("/polls", handlers.PostAPoll)
	router.POST("/polls/:pollID/:optionID", handlers.VoteOnAPoll)
	router.PATCH("/polls", handlers.EditAPoll)
	router.DELETE("/polls/:id", handlers.DeleteAPollByID)
	router.PATCH("/polls/:id", handlers.CloseVote)

	if err := router.Run("localhost:8080"); err != nil {
		log.Printf("Eroor running sercer: %v", err)
	}
}
