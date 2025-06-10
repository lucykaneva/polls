package models

import (
	"slices"
)

var polls []Poll

func GetAllPolls() *[]Poll {
	return &polls
}

// func EditAPollByID (id int, newPoll Poll){
// 	if (newPoll.Question == ""){
// 		newPoll.Question =
// 	}
// }

func AddAPoll(p Poll) {
	polls = append(polls, p)
}

func DeleteAPollByID(id string) bool {
	for i := 0; i < len(polls); i++ {
		if polls[i].ID == id {
			polls = slices.Delete(polls, i, i+1)
			return true
		}
	}
	return false
}
