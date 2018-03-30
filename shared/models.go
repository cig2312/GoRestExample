package shared

type Recipe struct {
	RecipeID   int
	Name       string
	PrepTime   string
	Difficulty int
	Vegetarian bool
}

type NewRecipe struct {
	Name       string
	PrepTime   string
	Difficulty int
	Vegetarian bool
}

type Rating struct {
	RecipeID int
	Ratings  []int
}

type RatingRequest struct {
	Rating int
}

type Token struct {
	Token string
}
