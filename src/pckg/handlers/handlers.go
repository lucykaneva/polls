package handlers

import (
	"log"
	"net/http"
	models "polls/src/pckg/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllPolls(c *gin.Context) {
	c.JSON(http.StatusOK, models.GetAllPolls())
}

func PostAPoll(c *gin.Context) {
	var newPoll models.Poll
	if err := c.BindJSON(&newPoll); err != nil {
		log.Print("Fail to decode")
		return
	}
	if models.IsThereADuplicateQuestion(newPoll.Question) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "poll's question is duplicated"})
		return
	}
	if models.IsThereADuplicateAnswerForAPoll(newPoll) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "poll's answer is duplicated"})
		return
	}
	newPoll.ID = uuid.New().String()
	newPoll.IsClosed = false
	for i := 0; i < len(newPoll.AnswerOptions); i++ {
		newPoll.AnswerOptions[i].ID = uuid.New().String()
	}
	models.AddAPoll(newPoll)
	c.JSON(http.StatusCreated, newPoll)
}

func GetAPollByID(c *gin.Context) {
	id := c.Param("id")
	if p, err := models.GetAPollByID(id); err == nil {
		c.JSON(http.StatusOK, *p)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "poll not found"})
}

func VoteOnAPoll(c *gin.Context) {
	idPoll := c.Param("pollID")
	idOption := c.Param("optionID")
	if p, err := models.GetAPollByID(idPoll); p.IsClosed && err == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "closed vote"})
		return
	}
	if _, err := models.GetAnOptionByIDs(idPoll, idOption); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "invalid id for options or for the poll"})
		return
	}
	o, _ := models.GetAnOptionByIDs(idPoll, idOption)
	o.Votes++
	p, _ := models.GetAPollByID(idPoll)
	c.JSON(http.StatusOK, *p)
}

// I removed the two functions for closing vote and editing and just united them in this new ModifyAPollByID on the single endpoint PATCH /polls/:id
func ModifyAPollByID(c *gin.Context) {
	idPoll := c.Param("id")

	var newUpdatedPoll models.Poll

	if err := c.BindJSON(&newUpdatedPoll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "update is unsuccessful, error in binding??"})
		return
	}
	if _, err := models.GetAPollByID(idPoll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "update is unsuccessful, invalid id??"})
		return
	}
	if models.IsThereADuplicateQuestion(newUpdatedPoll.Question) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "poll's question is duplicated"})
		return
	}
	if models.IsThereADuplicateAnswerForAPoll(newUpdatedPoll) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "there is a duplication in answers"})
		return
	}
	//setting new ids for the answer options and the setting the votes to 0
	for i := 0; i < len(newUpdatedPoll.AnswerOptions); i++ {
		newUpdatedPoll.AnswerOptions[i].ID = uuid.NewString()
		newUpdatedPoll.AnswerOptions[i].Votes = 0
	}
	p, _ := models.GetAPollByID(idPoll)
	p.EditAPollByID(newUpdatedPoll)
	c.JSON(http.StatusCreated, *p)
}

func DeleteAPollByID(c *gin.Context) {
	id := c.Param("id")
	if !models.DeleteAPollByID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "delete is unsuccessful, invalid id??"})
		return
	}
	c.JSON(http.StatusOK, models.GetAllPolls())
}
