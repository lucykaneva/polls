package models

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func ValidateUser(u User) (bool, error) {

	currentUsers, err := GetAllUsers()
	if err != nil {
		return false, err
	}

	for _, cu := range currentUsers {
		if cu.Username == u.Username {
			err = bcrypt.CompareHashAndPassword([]byte(cu.Password), []byte(u.Password))
			if err != nil {
				return false, nil
			} else {
				return true, nil
			}
		}
	}
	return false, nil

}
func IsUserOkayToGoToDatabase(u User) bool {

	currentUsers, err := GetAllUsers()
	if err != nil {
		log.Println(err)
		return false
	}
	if u.Password == "" {
		return false
	}
	for _, cu := range currentUsers {
		if cu.Username == u.Username {
			return false
		}
	}
	return true

}
