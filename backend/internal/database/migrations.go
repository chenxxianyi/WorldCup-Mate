package database

import (
	"log"

	"worldcup-mate/internal/models"
)

func EnsureUserUsernameIsNotUnique() error {
	if !DB.Migrator().HasTable(&models.User{}) {
		return nil
	}
	if !DB.Migrator().HasIndex(&models.User{}, "idx_users_username") {
		return nil
	}
	if err := DB.Migrator().DropIndex(&models.User{}, "idx_users_username"); err != nil {
		return err
	}
	log.Println("dropped unique index users.idx_users_username")
	return nil
}
