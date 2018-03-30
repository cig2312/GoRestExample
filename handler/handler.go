package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	. "gorestexample/config"
	. "gorestexample/shared"
)

type RecipeHandler struct {
	Conf *Configuration
}

func newRecipeHandler(conf *Configuration) *RecipeHandler {
	return &RecipeHandler{Conf: conf}
}

func RegisterHandlers(router *mux.Router, conf *Configuration) {

	rh := newRecipeHandler(conf)

	router.HandleFunc(LOGIN, rh.LoginHandler).Methods("GET")
	router.HandleFunc(RECIPES_ROUTE, rh.RecipesHandler).Methods("GET", "POST")
	router.HandleFunc(RECIPE_ROUTE, rh.RecipeHandler).Methods("GET", "DELETE", "PUT")
	router.HandleFunc(RECIPE_RATING_ROUTE, rh.RecipeRatingHandler).Methods("POST")
	router.HandleFunc(RECIPE_SEARCH_ROUTE, rh.SearchHandler).Methods("GET")

}

func (rh *RecipeHandler) LoginHandler(res http.ResponseWriter, req *http.Request) {

	conf := rh.Conf
	log := conf.AppLogger

	log.Logger.Printf("\n Route :%s", LOGIN)

	//validate user credentials by the credentails passed by Basic Auth

	username, password, _ := req.BasicAuth()

	if username != "admin" && password != "admin" {
		log.Logger.Println(INCORRECT_USERNAME_PASSWORD)
		processError(res, errors.New(INCORRECT_USERNAME_PASSWORD))
		return
	}

	//creating a claim

	claims := jwt.MapClaims{
		"username":  username,
		"ExpiresAt": 20000,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(conf.AppKeys.PrivateKey)

	if err != nil {
		log.Logger.Println("Failed to create Token")
		processError(res, errors.New(INTERNAL_SERVER_ERROR))
		return
	}

	processJson(res, http.StatusOK, map[string]string{"Token": tokenString})

}

func (rh *RecipeHandler) RecipesHandler(res http.ResponseWriter, req *http.Request) {

	conf := rh.Conf
	log := conf.AppLogger

	switch req.Method {

	case "GET":

		log.Logger.Printf("\n GET::Route : %s", RECIPES_ROUTE)
		recipes, err := conf.Db.GetAllRecipes()

		if err != nil {
			log.Logger.Println("error::", err)
			processError(res, err)
			return
		}

		processJson(res, http.StatusOK, recipes)

	case "POST":

		log.Logger.Printf("\n POST::Route : %s", RECIPES_ROUTE)

		if !isTokenValid(req, conf.AppKeys.PublicKey) {
			log.Logger.Println(INVALID_TOKEN)
			processError(res, errors.New(INVALID_TOKEN))
			return
		}

		recipeRequest, err := ValidateAndDecodeRecipeRequest(req)

		if err != nil {
			log.Logger.Println("error::", BAD_REQUEST)
			processError(res, errors.New(BAD_REQUEST))
			return
		}

		err = conf.Db.CreateRecipe(recipeRequest)

		if err != nil {
			log.Logger.Println("error::", err)
			processError(res, err)
			return
		}

		processJson(res, http.StatusCreated, map[string]string{"Result": "Successfully created"})

	}
}

func (rh *RecipeHandler) RecipeHandler(res http.ResponseWriter, req *http.Request) {

	conf := rh.Conf
	log := conf.AppLogger

	pathParam := mux.Vars(req)
	recipeId := pathParam["id"]

	switch req.Method {

	case "GET":

		log.Logger.Printf("\n GET::Route : %s :: Pathparam::%s", RECIPE_ROUTE, recipeId)

		recipe, err := conf.Db.GetRecipeById(recipeId)
		if err != nil {
			log.Logger.Println("error::", err)
			processError(res, err)
			return
		}
		processJson(res, http.StatusOK, recipe)

	case "DELETE":

		log.Logger.Printf("\n DELETE::Route : %s :: Pathparam::%s", RECIPE_ROUTE, recipeId)

		if !isTokenValid(req, conf.AppKeys.PublicKey) {
			log.Logger.Println("error::", INVALID_TOKEN)
			processError(res, errors.New(INVALID_TOKEN))
			return
		}
		if !isRecipePresent(recipeId, conf) {
			log.Logger.Println("error::", ID_DOES_NOT_EXIST)
			processError(res, errors.New(ID_DOES_NOT_EXIST))
			return
		}

		err := conf.Db.DeleteRecipeById(recipeId)

		if err != nil {
			log.Logger.Println("error::", err)
			processError(res, err)
			return
		}

		processJson(res, http.StatusOK, map[string]string{"result": "Recipe successfully deleted"})

	case "PUT":

		log.Logger.Printf("\n PUT::Route : %s :: Pathparam::%s", RECIPE_ROUTE, recipeId)

		if !isTokenValid(req, conf.AppKeys.PublicKey) {
			log.Logger.Println("error::", INVALID_TOKEN)
			processError(res, errors.New(INVALID_TOKEN))
			return
		}

		recipeRequest, err := ValidateAndDecodeRecipeRequest(req)

		if err != nil {
			log.Logger.Println("error::", BAD_REQUEST)
			processError(res, errors.New(BAD_REQUEST))
			return
		}
		if !isRecipePresent(recipeId, conf) {
			log.Logger.Println("error::", ID_DOES_NOT_EXIST)
			processError(res, errors.New(ID_DOES_NOT_EXIST))
			return
		}

		err = conf.Db.ModifyRecipe(recipeRequest, recipeId)
		if err != nil {
			log.Logger.Println("error::", err)
			processError(res, err)
			return
		}

		processJson(res, http.StatusOK, map[string]string{"result": "Recipe successfully modfied"})

	}
}

func (rh *RecipeHandler) RecipeRatingHandler(res http.ResponseWriter, req *http.Request) {

	conf := rh.Conf
	log := conf.AppLogger

	pathParam := mux.Vars(req)
	recipeID := pathParam["id"]

	log.Logger.Printf("\n POST::Route : %s :: Pathparam::%s", RECIPE_RATING_ROUTE, recipeID)

	if !isTokenValid(req, conf.AppKeys.PublicKey) {
		log.Logger.Println("error::", INVALID_TOKEN)
		processError(res, errors.New(INVALID_TOKEN))
		return
	}

	decoder := json.NewDecoder(req.Body)

	var ratingRequest RatingRequest

	err := decoder.Decode(&ratingRequest)

	if err != nil {
		log.Logger.Println("error::", err)
		processError(res, errors.New(INTERNAL_SERVER_ERROR))
		return
	}

	if !isRecipePresent(recipeID, conf) {
		log.Logger.Println("error::", ID_DOES_NOT_EXIST)
		processError(res, errors.New(ID_DOES_NOT_EXIST))
		return
	}

	if ratingRequest.Rating > 5 || ratingRequest.Rating < 1 {
		log.Logger.Println("error::Bad request: rating has to be between 1-5")
		processError(res, errors.New(BAD_REQUEST))
		return
	}

	err = conf.Db.RateRecipe(ratingRequest.Rating, recipeID)

	if err != nil {
		log.Logger.Println("error::", err)
		processError(res, err)
		return
	}

	processJson(res, http.StatusOK, map[string]string{"result": "Recipe Ratings added"})

}

func (rh *RecipeHandler) SearchHandler(res http.ResponseWriter, req *http.Request) {

	conf := rh.Conf
	log := conf.AppLogger

	param := req.FormValue("search")

	log.Logger.Printf("\n GET::Route : %s :: Queryparam::%s", RECIPE_SEARCH_ROUTE, param)

	recipes, err := conf.Db.Search(param)

	if err != nil {
		log.Logger.Println("error::", err)
		processError(res, err)
		return
	}

	processJson(res, http.StatusOK, recipes)

}
