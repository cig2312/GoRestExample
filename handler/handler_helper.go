package handler

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/pkg/errors"

	. "gorestexample/config"
	. "gorestexample/shared"
)

func processJson(res http.ResponseWriter, httpCode int, data interface{}) {

	response, _ := json.Marshal(data)
	res.WriteHeader(httpCode)
	res.Write(response)
}

func processError(res http.ResponseWriter, err error) {

	errString := err.Error()

	switch errString {

	case "Internal Server Error: sql: no rows in result set":
		processJson(res, http.StatusNotFound, map[string]string{"result": "Not Found"})

	case INCORRECT_USERNAME_PASSWORD:
		processJson(res, http.StatusUnauthorized, map[string]string{"result": INCORRECT_USERNAME_PASSWORD})

	case INVALID_TOKEN:
		processJson(res, http.StatusUnauthorized, map[string]string{"result": INVALID_TOKEN})

	case ID_ALREADY_PRESENT:
		processJson(res, http.StatusBadRequest, map[string]string{"result": ID_ALREADY_PRESENT})

	case ID_DOES_NOT_EXIST:
		processJson(res, http.StatusNotFound, map[string]string{"result": ID_DOES_NOT_EXIST})

	case BAD_REQUEST:
		processJson(res, http.StatusNotFound, map[string]string{"result": BAD_REQUEST})

	default:
		processJson(res, http.StatusInternalServerError, map[string]string{"result": INTERNAL_SERVER_ERROR})
	}

}

func isRecipePresent(id string, conf *Configuration) bool {

	_, err := conf.Db.GetRecipeById(id)

	if err != nil {
		return false
	}

	return true
}

func isTokenValid(req *http.Request, publicKey *rsa.PublicKey) bool {

	token, _ := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if token != nil {
		if token.Valid {
			return true
		}
	}
	return false
}

func ValidateAndDecodeRecipeRequest(req *http.Request) (NewRecipe, error) {

	var recipe NewRecipe

	err := json.NewDecoder(req.Body).Decode(&recipe)
	if err != nil {
		return NewRecipe{}, errors.New("Invalid request body")
	}

	if recipe.Name == "" {
		return NewRecipe{}, errors.New("Name Cannot be empty")
	} else if recipe.Difficulty > 5 || recipe.Difficulty < 1 {
		return NewRecipe{}, errors.New("Difficulty has to be from 1-5")
	} else if recipe.PrepTime == "" {
		return NewRecipe{}, errors.New("Prep Time Name Cannot be empty")
	}
	return recipe, nil
}
