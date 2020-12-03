package service

import "fmt"

// DB db
type DB interface {
	Get(string) (string, error)
}

// GetFromDB get from db
func GetFromDB(db DB, key string) (string, error) {
	data, err := db.Get(key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("DB: %s", data), nil
}
