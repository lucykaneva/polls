package handlers

import (
	"log"
	"net/http"
	models "polls/src/pckg/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllPolls(c *gin.Context) {
	polls, _ := models.GetAllPolls()
	c.JSON(http.StatusOK, polls)
}

func PostAPoll(c *gin.Context) {

	var newPoll models.Poll

	err := c.ShouldBindJSON(&newPoll)
	if err != nil {
		c.JSON(501, gin.H{"message": "invalid request format"})
	}

	// var newPoll models.Poll
	// log.Println("just here")
	// if err := c.BindJSON(&newPoll); err != nil {
	// 	log.Println("here1")

	// 	log.Print("Fail to decode")
	// 	return
	// }
	if models.IsThereADuplicateQuestion(newPoll.Question) {
		log.Println("here2")
		c.JSON(http.StatusBadRequest, gin.H{"message": "poll's question is duplicated"})
		return
	}
	if models.IsThereADuplicateAnswerForAPoll(newPoll) {
		log.Println("here3")
		c.JSON(http.StatusBadRequest, gin.H{"message": "poll's answer is duplicated"})
		return
	}
	log.Println("here4")

	newPoll.IsClosed = false

	for i := 0; i < len(newPoll.AnswerOptions); i++ {
		newPoll.AnswerOptions[i].ID = primitive.NewObjectID()
		for j := 0; j < len(newPoll.AnswerOptions[i].Votes); j++ {
			newPoll.AnswerOptions[i].Votes[j].ID = primitive.NewObjectID()
		}
	}
	models.AddAPoll(newPoll)
	c.JSON(http.StatusCreated, newPoll)
}

// func GetAPollByID(c *gin.Context) {
// 	id := c.Param("id")
// 	idObjectPrim, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "delete is unsuccessful, parsing id to object primitive"})
// 		return
// 	}
// 	if p, err := models.GetAPollByID(idObjectPrim); err == nil {
// 		c.JSON(http.StatusOK, p)
// 		return
// 	}
// 	c.JSON(http.StatusNotFound, gin.H{"message": "poll not found"})
// }

// func VoteOnAPoll(c *gin.Context) {
// 	idPoll := c.Param("pollID")
// 	idOption := c.Param("optionID")
// 	idPollObjectPrim, err := primitive.ObjectIDFromHex(idPoll)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "delete is unsuccessful, parsing id to object primitive"})
// 		return
// 	}
// 	idOptionObjectPrim, _ := primitive.ObjectIDFromHex(idOption)
// 	if p, err := models.GetAPollByID(idPollObjectPrim); p.IsClosed && err == nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "closed vote"})
// 		return
// 	}
// 	if _, err := models.GetAnOptionByIDs(idPollObjectPrim, idOptionObjectPrim); err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "invalid id for options or for the poll"})
// 		return
// 	}
// 	o, _ := models.GetAnOptionByIDs(idPollObjectPrim, idOptionObjectPrim)
// 	//I don't have a userID yet
// 	currentVote := models.Vote{ID: primitive.NewObjectID(), UserID: "abc"}
// 	o.Votes = append(o.Votes, currentVote)
// 	p, _ := models.GetAPollByID(idPollObjectPrim)
// 	c.JSON(http.StatusOK, p)
// }

// I removed the two functions for closing vote and editing and just united them in this new ModifyAPollByID on the single endpoint PATCH /polls/:id
// func ModifyAPollByID(c *gin.Context) {
// 	idPoll := c.Param("id")
// 	idObjectPrim, err := primitive.ObjectIDFromHex(idPoll)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "delete is unsuccessful, parsing id to object primitive"})
// 		return
// 	}
// 	var newUpdatedPoll models.Poll

// 	if err := c.BindJSON(&newUpdatedPoll); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "update is unsuccessful, error in binding??"})
// 		return
// 	}
// 	if _, err := models.GetAPollByID(idObjectPrim); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "update is unsuccessful, invalid id??"})
// 		return
// 	}
// 	if models.IsThereADuplicateQuestion(newUpdatedPoll.Question) {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "poll's question is duplicated"})
// 		return
// 	}
// 	if models.IsThereADuplicateAnswerForAPoll(newUpdatedPoll) {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "there is a duplication in answers"})
// 		return
// 	}
// 	//setting new ids for the answer options and the setting the votes to 0
// 	for i := 0; i < len(newUpdatedPoll.AnswerOptions); i++ {
// 		newUpdatedPoll.AnswerOptions[i].ID = primitive.NewObjectID()
// 		newUpdatedPoll.AnswerOptions[i].Votes = []models.Vote{}
// 	}
// 	p, _ := models.GetAPollByID(idObjectPrim)
// 	p.EditAPollByID()
// 	c.JSON(http.StatusCreated, p)
// }

// func DeleteAPollByID(c *gin.Context) {
// 	id := c.Param("id")
// 	idObjectPrim, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "delete is unsuccessful, parsing id to object primitive"})
// 		return
// 	}
// 	if !models.DeleteAPollByID(idObjectPrim) {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "delete is unsuccessful, invalid id??"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, models.GetAllPolls())
// }
