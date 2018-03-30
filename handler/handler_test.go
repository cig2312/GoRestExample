package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "gorestexample/config"
	. "gorestexample/handler"
	. "gorestexample/shared"
)

func openFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0666)
}

func setupLogger() *log.Logger {
	// create file if not present
	file, err := os.OpenFile("RecipesAppTest.log", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Cannot open log file.")
		os.Exit(1)
	}
	Logger := log.New(file,
		"LOG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	return Logger
}

func setupMockKeys() AppKeys {

	privateKeyBytes, err := ioutil.ReadFile("../keys/app.rsa")
	if err != nil {
		log.Fatal(PRIVATE_KEY_FILE_ERROR)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		log.Fatal(PRIVATE_KEY_ERROR)
	}
	publicKeyBytes, err := ioutil.ReadFile("../keys/app.rsa.pub")
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

var recipesratings = []Rating{
	{
		RecipeID: 1,
		Ratings:  []int{2, 5, 3},
	},
	{
		RecipeID: 5,
		Ratings:  []int{2, 5, 3},
	},
}

var recipesOutput = []Recipe{
	{
		RecipeID:   1,
		Name:       "roti",
		PrepTime:   "10 mins ",
		Difficulty: 2,
		Vegetarian: true,
	},
	{
		RecipeID:   5,
		Name:       "roticurry",
		PrepTime:   "10 mins ",
		Difficulty: 3,
		Vegetarian: true,
	},
}

var recipe = Recipe{

	RecipeID:   1,
	Name:       "roti",
	PrepTime:   "10 mins ",
	Difficulty: 2,
	Vegetarian: true,
}

var recipeString = []byte(`{
	RecipeID:   1,
	Name:       "roti",
	PrepTime:   "10 mins ",
	Difficulty: 2,
	Vegetarian: true,
}`)

var updatedRecipeString = []byte(`{
	RecipeID:   1,
	Name:       "steak",
	PrepTime:   "50 mins ",
	Difficulty: 5,
	Vegetarian: true,
}`)

var mockKeys = setupMockKeys()

var logger = AppLogger{Logger: setupLogger()}

type mockDB struct{}

var mockdb = &mockDB{}

var keys = AppKeys{
	PublicKey:  mockKeys.PublicKey,
	PrivateKey: mockKeys.PrivateKey,
}

var conf = &Configuration{mockdb, logger, keys}

var handler = &RecipeHandler{Conf: conf}

func getToken() string {
	response := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	req.SetBasicAuth("admin", "admin")
	http.HandlerFunc(handler.LoginHandler).ServeHTTP(response, req)
	var token Token
	json.Unmarshal(response.Body.Bytes(), &token)

	return token.Token
}

func (mockdb *mockDB) GetAllRecipes() ([]Recipe, error) {
	return recipesOutput, nil
}

func (mockdb *mockDB) GetRecipeById(id string) (Recipe, error) {
	for _, recipe := range recipesOutput {
		id, _ := strconv.Atoi(id)
		if recipe.RecipeID == id {
			return recipe, nil
		}
	}
	return Recipe{}, errors.New("recipe not found")
}

func (mockdb *mockDB) CreateRecipe(newrecipe NewRecipe) error {

	recipe := Recipe{
		RecipeID:   3,
		Name:       newrecipe.Name,
		Difficulty: newrecipe.Difficulty,
		PrepTime:   newrecipe.PrepTime,
		Vegetarian: newrecipe.Vegetarian,
	}

	recipesOutput = append(recipesOutput, recipe)
	return nil
}

func (mockdb *mockDB) DeleteRecipeById(id string) error {
	for index, recipe := range recipesOutput {
		id, _ := strconv.Atoi(id)
		if recipe.RecipeID == id {
			recipesOutput = append(recipesOutput[:index], recipesOutput[index+1:]...)
			return nil
		}
	}
	return errors.New("recipe not found")
}

func (mockdb *mockDB) ModifyRecipe(recipe NewRecipe, id string) error {
	integerID, _ := strconv.Atoi(id)
	for _, recipeValue := range recipesOutput {
		if recipeValue.RecipeID == integerID {
			recipeValue.Name = recipe.Name
			recipeValue.Difficulty = recipe.Difficulty
			recipeValue.PrepTime = recipe.PrepTime
			recipeValue.Vegetarian = recipe.Vegetarian
			return nil
		}
	}
	return errors.New("recipe not found")
}

func (mockdb *mockDB) Search(search string) ([]Recipe, error) {
	var searchedRecipes []Recipe
	for _, recipe := range recipesOutput {
		if strings.Contains(recipe.Name, search) {
			searchedRecipes = append(searchedRecipes, recipe)
		}
		fmt.Print("ASDAOSDASDAIIIIIIWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWW")
		fmt.Print(searchedRecipes)
		return searchedRecipes, nil
	}

	return nil, errors.New("recipe not found")
}

func (mockdb *mockDB) RateRecipe(rating int, id string) error {
	integerID, _ := strconv.Atoi(id)
	for _, recipeValue := range recipesratings {
		if recipeValue.RecipeID == integerID {
			recipeValue.Ratings = append(recipeValue.Ratings, rating)
			return nil
		}
	}
	return errors.New("recipe not found")
}

var _ = Describe("Http Handler", func() {

	Context("Get all recipes :: GET :: /recipes ", func() {

		It("returns a list of all the recipes", func() {

			response := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/recipes", nil)
			http.HandlerFunc(handler.RecipesHandler).ServeHTTP(response, req)

			var recipes []Recipe
			json.Unmarshal(response.Body.Bytes(), &recipes)

			Expect(recipes).To(Equal(recipesOutput))
			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Get a Recipe by recipe ID :: GET :: /recipes/{id}", func() {

		It("returns a recipe based on the recipe ID Given in the path param", func() {

			router := mux.NewRouter()
			router.HandleFunc(RECIPE_ROUTE, handler.RecipeHandler)

			response := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/recipes/5", nil)

			router.ServeHTTP(response, req)

			var recipe Recipe
			json.Unmarshal(response.Body.Bytes(), &recipe)

			Expect(recipe).To(Equal(recipesOutput[1]))
			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Create Recipe :: POST :: /recipes ", func() {

		It("creates a recipe based on the given request body", func() {

			router := mux.NewRouter()
			router.HandleFunc(RECIPES_ROUTE, handler.RecipeHandler)

			payload := []byte(`
				{
				"Name":       "roti",
				"PrepTime":   "10 mins ",
				"Difficulty": 2,
				"Vegetarian": true
				}
				`)
			req, _ := http.NewRequest("POST", "/recipes", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			token := getToken()

			req.Header.Set("authorization", "bearer "+token)

			response := httptest.NewRecorder()

			router.ServeHTTP(response, req)
			fmt.Println(response)
			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Update Recipe :: PUT :: /recipes/{id} ", func() {

		It("update recipe based on path param and request body", func() {

			router := mux.NewRouter()
			router.HandleFunc(RECIPE_ROUTE, handler.RecipeHandler)

			payload := []byte(`
				{
				"Name":       "roti",
				"PrepTime":   "10 mins ",
				"Difficulty": 2,
				"Vegetarian": true
				}
				`)
			req, _ := http.NewRequest("PUT", "/recipes/1", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			token := getToken()
			req.Header.Set("authorization", "bearer "+token)

			response := httptest.NewRecorder()

			router.ServeHTTP(response, req)

			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Delete Recipe :: DELETE :: /recipes/2 ", func() {

		It("deletes a recipe based on the given id in path param", func() {

			router := mux.NewRouter()
			router.HandleFunc(RECIPE_ROUTE, handler.RecipeHandler)

			req, _ := http.NewRequest("DELETE", "/recipes/1", nil)
			req.Header.Set("Content-Type", "application/json")

			token := getToken()
			req.Header.Set("authorization", "bearer "+token)

			response := httptest.NewRecorder()

			router.ServeHTTP(response, req)

			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Rate Recipe :: POST :: /recipes ", func() {

		It("rates a recipe based on the given request body and path param", func() {

			router := mux.NewRouter()
			router.HandleFunc(RECIPE_RATING_ROUTE, handler.RecipeRatingHandler)

			payload := []byte(`
				{
				"Rating":   1
				}
				`)
			req, _ := http.NewRequest("POST", "/recipes/5/rating", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			token := getToken()
			req.Header.Set("authorization", "bearer "+token)

			response := httptest.NewRecorder()

			router.ServeHTTP(response, req)

			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Search Recipe :: GET :: /recipes ", func() {

		It("searches for a recipe based on the param", func() {

			response := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "recipes/search/?search=roti", nil)
			http.HandlerFunc(handler.SearchHandler).ServeHTTP(response, req)

			var recipes []Recipe
			json.Unmarshal(response.Body.Bytes(), &recipes)

			Expect(recipes).To(Equal(recipesOutput))
			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})

})
