package handlers

import (
	"errors"
	"log"
	"net/http"
	models "polls/src/pckg/models"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var polls []models.Poll

func GetAllPolls(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, polls)
}

func CloseVote(c *gin.Context) {
	id := c.Param("id")
	if p, err := getAPollByID(id); err == nil {
		p.IsClosed = true
		c.IndentedJSON(http.StatusOK, *p)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "poll to be closed not found"})
}

func PostAPoll(c *gin.Context) {
	var newPoll models.Poll
	if err := c.BindJSON(&newPoll); err != nil {
		log.Print("Fail to decode")
		return
	}
	newPoll.ID = uuid.New().String()
	newPoll.IsClosed = false
	for i := 0; i < len(newPoll.AnswerOptions); i++ {
		newPoll.AnswerOptions[i].ID = uuid.New().String()
	}
	polls = append(polls, newPoll)
	c.IndentedJSON(http.StatusCreated, newPoll)
}

func GetAPollByID(c *gin.Context) {
	id := c.Param("id")
	if p, err := getAPollByID(id); err == nil {
		c.IndentedJSON(http.StatusOK, *p)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "poll not found"})
}

func VoteOnAPoll(c *gin.Context) {
	idPoll := c.Param("pollID")
	idOption := c.Param("optionID")
	if p,err:=getAPollByID(idPoll); p.IsClosed == true && err==nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "closed vote"})
		return
	}
	if o, err := getAnOptionByIDs(idPoll, idOption); err == nil {
		o.Votes++
		p, _ := getAPollByID(idPoll)
		c.IndentedJSON(http.StatusOK, *p)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "invalid id for options or for the poll"})
}

func EditAPoll(c *gin.Context) {
	var newUpdatedPoll models.Poll
	if err := c.BindJSON(&newUpdatedPoll); err != nil {
		return
	}
	if p, err := getAPollByID(newUpdatedPoll.ID); err == nil {
		*p = newUpdatedPoll
		c.IndentedJSON(http.StatusCreated, *p)
		return
	}
	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "update is unsuccessful, invalid id??"})
}

func DeleteAPollByID(c *gin.Context) {
	id := c.Param("id")
	for i := 0; i < len(polls); i++ {
		if polls[i].ID == id {
			polls = slices.Delete(polls, i, i+1)
			c.IndentedJSON(http.StatusCreated, polls)
			return
		}
	}
	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "delete is unsuccessful, invalid id??"})
}

func getAPollByID(id string) (*models.Poll, error) {
	for i := 0; i < len(polls); i++ {
		if polls[i].ID == id {
			return &polls[i], nil
		}
	}
	return nil, errors.New("poll not found")
}
func getAnOptionByIDs(idPoll string, idOption string) (*models.Option, error) {
	if p, err := getAPollByID(idPoll); err == nil {
		for i := 0; i < len(p.AnswerOptions); i++ {
			if p.AnswerOptions[i].ID == idOption {
				return &p.AnswerOptions[i], nil
			}
		}
	}

	return nil, errors.New("option not found")
}
