package config

import (
	"database/sql"
	"log"

	. "gorestexample/api"
	. "gorestexample/shared"
)

//Configuration structs holds all our components of the app
type Configuration struct {
	Db        RecipeStore
	AppLogger AppLogger
	AppKeys   AppKeys
}

// NewConfiguration injects the config struct with our instances of logger, db etc etc
func NewConfiguration(db *sql.DB, logger *log.Logger, keys AppKeys) *Configuration {

	RecipeDb := &RecipeDataBase{Db: db}
	appLog := AppLogger{Logger: logger}
	appKeys := keys
	Conf := Configuration{Db: RecipeDb, AppLogger: appLog, AppKeys: appKeys}
	return &Conf
}
