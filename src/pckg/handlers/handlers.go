package handlers

import (
	"log"
	"net/http"
	auth "polls/src/pckg/auth"
	models "polls/src/pckg/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SignUp(c *gin.Context) {
	var newUser models.User
	err := c.ShouldBindJSON(&newUser)
	if err != nil {
		c.JSON(501, gin.H{"message": "invalid request format"})
	}
	err = models.AddAUser(newUser)
	if err != nil {
		c.JSON(501, gin.H{"message": "couldn't add to database"})
	}
}

func LogIn(c *gin.Context) {
	var userToLogIn models.User
	err := c.ShouldBindJSON(&userToLogIn)
	if err != nil {
		c.JSON(501, gin.H{"message": "invalid request format"})
	}
	if b, err := models.ValidateUser(userToLogIn); err != nil || !b {
		log.Print(err, b)
		c.JSON(501, gin.H{"message": "invalid user"})
	} else {
		u, err := models.GetAUserByUsername(userToLogIn.Username)
		if err != nil {
			c.JSON(501, gin.H{"message": "can't find user"})
		}
		tokenString, err := auth.CreateToken(u.ID, u.Username)
		if err != nil {
			c.JSON(501, gin.H{"message": "can't create token"})
		}
		c.SetCookie("token", tokenString, 3600, "/", "", true, true)
		c.JSON(200, gin.H{"message": "login successful", "token": tokenString})
	}
}
func LogOut(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", true, true)
	c.JSON(200, gin.H{"message": "logout successful"})
}

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

	if models.IsThereADuplicateQuestion(newPoll.Question) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "poll's question is duplicated"})
		return
	}
	if models.IsThereADuplicateAnswerForAPoll(newPoll) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "poll's answer is duplicated"})
		return
	}

	newPoll.IsClosed = false
	userId, _ := c.Get("id")
	newPoll.CreatedUserID = userId.(primitive.ObjectID)

	newPoll.ID = primitive.NewObjectID()
	for i := 0; i < len(newPoll.AnswerOptions); i++ {
		newPoll.AnswerOptions[i].ID = primitive.NewObjectID()
		newPoll.AnswerOptions[i].Votes = []models.Vote{}
	}
	models.AddAPoll(newPoll)
	c.JSON(http.StatusCreated, newPoll)
}

func GetAPollByID(c *gin.Context) {
	id := c.Param("id")
	idObjectPrim, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "can't parse the parameter from string to primitive.ObjectID"})
		return
	}
	if p, err := models.GetAPollByID(idObjectPrim); err == nil {
		c.JSON(http.StatusOK, p)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "poll not found"})
}

func VoteOnAPoll(c *gin.Context) {
	idPoll := c.Param("pollID")
	idOption := c.Param("optionID")
	idPollObjectPrim, err := primitive.ObjectIDFromHex(idPoll)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "can't parse the parameter from string to primitive.ObjectID"})
		return
	}
	idOptionObjectPrim, _ := primitive.ObjectIDFromHex(idOption)
	if p, err := models.GetAPollByID(idPollObjectPrim); p.IsClosed && err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "closed vote"})
		return
	}
	if _, err := models.GetAnOptionByIDs(idPollObjectPrim, idOptionObjectPrim); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "invalid id for options or for the poll"})
		return
	}
	p, _ := models.GetAPollByID(idPollObjectPrim)

	userID, _ := c.Get("id")
	//creatin a new vote + appending it to a poll
	currentVote := models.Vote{ID: primitive.NewObjectID(), UserID: userID.(primitive.ObjectID)}
	for i := 0; i < len(p.AnswerOptions); i++ {
		if p.AnswerOptions[i].ID == idOptionObjectPrim {
			p.AnswerOptions[i].Votes = append(p.AnswerOptions[i].Votes, currentVote)
		}
	}

	//we need the new poll with the votes to again use the EditAPoll function that works with the database
	newPoll := p
	err = p.EditAPollByID(newPoll)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "couldn't update the poll with the vote in the db"})
		return
	}
	p, _ = models.GetAPollByID(idPollObjectPrim)
	c.JSON(http.StatusOK, p)
}

// this handles closing vote and editing a poll on the single endpoint PATCH /polls/:id
func ModifyAPollByID(c *gin.Context) {
	idPoll := c.Param("id")
	idObjectPrim, err := primitive.ObjectIDFromHex(idPoll)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "can't parse the parameter from string to primitive.ObjectID"})
		return
	}
	var newUpdatedPoll models.Poll

	if err := c.BindJSON(&newUpdatedPoll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "update is unsuccessful, error in binding"})
		return
	}
	if _, err := models.GetAPollByID(idObjectPrim); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "update is unsuccessful, invalid id"})
		return
	}
	if models.IsThereADuplicateAnswerForAPoll(newUpdatedPoll) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "there is a duplication in answer options"})
		return
	}
	p, _ := models.GetAPollByID(idObjectPrim)
	userRequestID, _ := c.Get("id")
	if p.CreatedUserID != userRequestID.(primitive.ObjectID) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "you are not the creator of this poll, so you can't edit it"})
		return
	}
	//setting new ids for the answer options and the setting the votes to 0
	for i := 0; i < len(newUpdatedPoll.AnswerOptions); i++ {
		newUpdatedPoll.AnswerOptions[i].ID = primitive.NewObjectID()
		newUpdatedPoll.AnswerOptions[i].Votes = []models.Vote{}
	}
	//newUpdatedPoll.ID = primitive.NewObjectID()

	log.Println("", newUpdatedPoll)
	log.Println("", p)
	err = p.EditAPollByID(newUpdatedPoll)
	if err != nil {
		log.Println("err editing a poll", err)
		c.JSON(500, gin.H{"error": "cant edit poll", "details": err.Error()})
		return
	}
	new, _ := models.GetAPollByID(p.ID)
	c.JSON(http.StatusOK, new)
}

func DeleteAPollByID(c *gin.Context) {
	id := c.Param("id")
	idObjectPrim, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "delete is unsuccessful, parsing id to object primitive"})
		return
	}
	if !models.DeleteAPollByID(idObjectPrim) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "delete is unsuccessful, invalid id"})
		return
	}
	polls, _ := models.GetAllPolls()
	c.JSON(http.StatusOK, polls)
}
