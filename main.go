package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	. "gorestexample/config"
	. "gorestexample/handler"
	. "gorestexample/shared"
)

/* This is the main class where all the components are initialzed before the program logic is executed
Components like logger, Database, Cache can be declared here. I have chosen to used notmal methods instead of init methods
becuase we have more control over normal methodos as ompared to init() methods */

// Initialize logger
func setupLogger() *log.Logger {

	// create file if not present
	file, err := os.OpenFile(LOG_FILE, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(CANNOT_OPEN_LOG_FILE)
		os.Exit(1)
	}
	Logger := log.New(file,
		"LOG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	return Logger
}

// create file if doesnt exist
func openFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0666)
}

//connect to DB
func connectDB() (*sql.DB, error) {

	var err error
	Db, err := sql.Open("postgres", "postgresql://gorest:gorest@postgres/gorest?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(SUCCESS_MESSAGE)
	}

	if err = Db.Ping(); err != nil {
		log.Fatal(err)
	}

	return Db, err
}

// read and generate public and private keys
func generateKeys() AppKeys {

	privateKeyBytes, err := ioutil.ReadFile(PathToPrivateKey)
	if err != nil {
		log.Fatal(PRIVATE_KEY_FILE_ERROR)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		log.Fatal(PRIVATE_KEY_ERROR)
	}
	publicKeyBytes, err := ioutil.ReadFile(PathToPublicKey)
	if err != nil {
		log.Fatalf(PUBLIC_KEY_FILE_ERROR)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		log.Fatalf(PUBLIC_KEY_ERROR)
	}

	return AppKeys{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

/* Once creation of all objects are complete inject these instances into config struct.
This struct can be used later in our program to access the functions of config instances
*/

func serveWeb() {

	db, err := connectDB()
	if err != nil {
		log.Fatal(DATABASE_CONNECTION_ERROR, err)
	}

	appKeys := generateKeys()
	logger := setupLogger()
	conf := NewConfiguration(db, logger, appKeys)
	router := mux.NewRouter()

	RegisterHandlers(router, conf)

	log.Fatal(http.ListenAndServe(":5000", router))

}

func main() {
	serveWeb()
}
