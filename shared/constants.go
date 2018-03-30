package shared

//Routes
const (
	RECIPES_ROUTE       = "/recipes"
	RECIPE_ROUTE        = "/recipes/{id}"
	RECIPE_RATING_ROUTE = "/recipes/{id}/rating"
	RECIPE_SEARCH_ROUTE = "/recipes/search/"
	LOGIN               = "/login"

	//Queries

	GET_ALL_RECIPE_QUERY   = "SELECT * FROM recipes"
	GET_RECIPE_BY_ID_QUERY = "SELECT * FROM recipes WHERE RecipeID = $1"
	CREATE_RECIPE_QUERY    = "INSERT INTO recipes(Name, PrepTime, Difficulty, Vegetarian) VALUES($1, $2, $3, $4) returning RecipeID"
	MODIFY_RECIPE_QUERY    = "UPDATE recipes SET Name = $1, Difficulty= $2, PrepTime= $3, Vegetarian=$4 WHERE RecipeID = $5"
	RATE_RECIPE_QUERY      = "UPDATE reciperatings SET rating = array_append(reciperatings.rating, $1 ) WHERE RecipeID = $2"
	SEARCH_QUERY           = "SELECT * FROM recipes WHERE Name LIKE '%' || $1 || '%'"
	DELETE_QUERY           = "DELETE FROM recipes * WHERE RecipeID = $1"
	RATINGS_INSERT_QUERY   = "INSERT INTO reciperatings VALUES($1)"

	//Messages

	INTERNAL_SERVER_ERROR       = "Internal Server Error"
	INCORRECT_USERNAME_PASSWORD = "Incorrect username or password"
	INVALID_TOKEN               = "Unauthorized. Invalid token"
	ID_ALREADY_PRESENT          = "ID already present"
	ID_DOES_NOT_EXIST           = "Id does not exist"
	BAD_REQUEST                 = "Invalid Input"

	CANNOT_OPEN_LOG_FILE      = "Cannot open log file."
	SUCCESS_MESSAGE           = "Server Up. Listening.... \n"
	PRIVATE_KEY_FILE_ERROR    = "Error reading private key from file"
	PRIVATE_KEY_ERROR         = "Error parsing private key"
	PUBLIC_KEY_FILE_ERROR     = "Error reading public key from file"
	PUBLIC_KEY_ERROR          = "Error parsing public key"
	DATABASE_CONNECTION_ERROR = "Not able to connect database."

	//Security Keys

	PathToPrivateKey = "./keys/app.rsa"
	PathToPublicKey  = "./keys/app.rsa.pub"

	//Files
	LOG_FILE = "RecipesApp.log"
)
