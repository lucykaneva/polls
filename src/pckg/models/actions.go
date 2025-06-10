package models

import (
	"errors"
)

func GetAPollByID(id string) (*Poll, error) {
	polls := GetAllPolls()
	for i := 0; i < len(*polls); i++ {
		if (*polls)[i].ID == id {
			return &(*polls)[i], nil
		}
	}
	return nil, errors.New("poll not found")
}
func GetAnOptionByIDs(idPoll string, idOption string) (*Option, error) {
	if p, err := GetAPollByID(idPoll); err == nil {
		for i := 0; i < len(p.AnswerOptions); i++ {
			if p.AnswerOptions[i].ID == idOption {
				return &p.AnswerOptions[i], nil
			}
		}
	}

	return nil, errors.New("option not found")
}
func IsThereADuplicateQuestion(question string) bool {
	polls := GetAllPolls()
	for _, p := range *polls {
		if p.Question == question {
			return true
		}
	}
	return false
}
func IsThereADuplicateAnswerForAPoll(p Poll) bool {
	answers := p.AnswerOptions
	seen := make(map[string]bool)
	for _, a := range answers {
		if seen[a.Content] {
			return true
		}
		seen[a.Content] = true
	}
	return false
}
