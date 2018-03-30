package api

import (
	"database/sql"

	"github.com/pkg/errors"

	. "gorestexample/shared"
)

type RecipeStore interface {
	GetAllRecipes() ([]Recipe, error)
	GetRecipeById(id string) (Recipe, error)
	CreateRecipe(recipe NewRecipe) error
	ModifyRecipe(recipe NewRecipe, recipeId string) error
	RateRecipe(rating int, id string) error
	Search(search string) ([]Recipe, error)
	DeleteRecipeById(id string) error
}

type RecipeDataBase struct {
	Db *sql.DB
}

func (Dbase *RecipeDataBase) GetAllRecipes() ([]Recipe, error) {

	Db := Dbase.Db

	rows, err := Db.Query(GET_ALL_RECIPE_QUERY)

	if err != nil {
		return nil, errors.New(INTERNAL_SERVER_ERROR)
	}
	defer rows.Close()

	recipes := make([]Recipe, 0)

	for rows.Next() {
		var recipe Recipe
		err := rows.Scan(&recipe.RecipeID, &recipe.Name, &recipe.PrepTime, &recipe.Difficulty, &recipe.Vegetarian)
		if err != nil {
			return nil, errors.New(INTERNAL_SERVER_ERROR)
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (Dbase *RecipeDataBase) GetRecipeById(id string) (Recipe, error) {

	Db := Dbase.Db

	row := Db.QueryRow(GET_RECIPE_BY_ID_QUERY, id)

	var recipe Recipe
	err := row.Scan(&recipe.RecipeID, &recipe.Name, &recipe.PrepTime, &recipe.Difficulty, &recipe.Vegetarian)
	if err != nil {
		return Recipe{}, errors.New(INTERNAL_SERVER_ERROR)
	}

	return recipe, nil
}

func (Dbase *RecipeDataBase) CreateRecipe(recipe NewRecipe) error {

	Db := Dbase.Db

	row := Db.QueryRow(CREATE_RECIPE_QUERY, recipe.Name, recipe.PrepTime, recipe.Difficulty, recipe.Vegetarian)

	var recipeID int
	row.Scan(&recipeID)

	_, err := Db.Exec(RATINGS_INSERT_QUERY, recipeID)

	if err != nil {
		return errors.New(INTERNAL_SERVER_ERROR)
	}

	return nil
}

func (Dbase *RecipeDataBase) ModifyRecipe(recipe NewRecipe, recipeId string) error {

	Db := Dbase.Db

	_, err := Db.Exec(MODIFY_RECIPE_QUERY, recipe.Name, recipe.Difficulty, recipe.PrepTime, recipe.Vegetarian, recipeId)

	if err != nil {
		return errors.New(INTERNAL_SERVER_ERROR)
	}

	return nil
}

func (Dbase *RecipeDataBase) RateRecipe(rating int, id string) error {

	Db := Dbase.Db

	_, err := Db.Exec(RATE_RECIPE_QUERY, rating, id)

	if err != nil {
		return errors.New(INTERNAL_SERVER_ERROR)
	}

	return nil
}

func (Dbase *RecipeDataBase) Search(search string) ([]Recipe, error) {

	Db := Dbase.Db

	rows, err := Db.Query(SEARCH_QUERY, search)

	if err != nil {
		return []Recipe{}, errors.New(INTERNAL_SERVER_ERROR)
	}
	defer rows.Close()

	recipes := make([]Recipe, 0)
	for rows.Next() {

		var recipe Recipe
		err := rows.Scan(&recipe.RecipeID, &recipe.Name, &recipe.PrepTime, &recipe.Difficulty, &recipe.Vegetarian)

		if err != nil {
			return nil, errors.New(INTERNAL_SERVER_ERROR)
		}

		recipes = append(recipes, recipe)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New(INTERNAL_SERVER_ERROR)
	}

	return recipes, nil

}

func (Dbase *RecipeDataBase) DeleteRecipeById(id string) error {

	Db := Dbase.Db

	_, err := Db.Query(DELETE_QUERY, id)

	if err != nil {
		return errors.New(INTERNAL_SERVER_ERROR)
	}

	return nil
}
