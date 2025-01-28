package database

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3" // Import du driver SQLite
)

// Database représente l'instance de la base de données
type Database struct {
	DB *sql.DB
}

var (
	dbInstance *Database
	once       sync.Once
)

// GetDatabaseInstance retourne l'instance unique de la base de données
func GetDatabaseInstance() *Database {
	once.Do(func() {
		var err error
		dbInstance = &Database{}
		dbInstance.DB, err = sql.Open("sqlite3", "./forum.db")
		if err != nil {
			panic("Erreur lors de l'ouverture de la base de données: " + err.Error())
		}
	})

	return dbInstance
}

func CloseDB() error {
	if dbInstance != nil && dbInstance.DB != nil {
		return dbInstance.DB.Close()
	}
	return nil
}
