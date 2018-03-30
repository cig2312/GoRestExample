#GoRestExample

### Prerequisites

Go 1.8 or higher
Docker 
Docker-compose 

### Running the project

Todo 
```
go run main.go
```
The project will now be running and listening to requests on port 5000


#### API definitions and Usage

| Name   | Method      | URL                    | Protected |   
| ---    | ---         | ---                    | ---       |
| Login  | `GET`       | `/login`               | ✘         |
| List   | `GET`       | `/recipes`             | ✘         |
| Create | `POST`      | `/recipes`             | ✓         |
| Get    | `GET`       | `/recipes/{id}`        | ✘         |
| Update | `PUT`       | `/recipes/{id}`        | ✓         |
| Delete | `DELETE`    | `/recipes/{id}`        | ✓         |
| Rate   | `POST`      | `/recipes/{id}/rating` | ✘         |
| Search | `GET`       | `/recipes/search?q=""` | ✘         |

Examples: 

1. Login : Login to get authorization token for accessing protected endpoints
           Use admin/admin as username/password for the "Basic Authorization" method  
           URL: `/login`
           Method Type: `GET` 
           Headers: `authorization: Basic YWRtaW46YWRtaW4=`
           Sample Response: `{"Token":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjIwMDAwLCJ1c2VybmFtZSI6ImFkbWluIn0"}`

2. List all Recipes:   
           URL: `/recipes`
           Method Type: `GET` 
           Headers: None 
           Sample Response: 
            `[{"RecipeID":1,"Name":"roti","PrepTime":"10 mins ","Difficulty":2,"Vegetarian":true},{"RecipeID":6,"Name":"butter roti","PrepTime":"10 mins ","Difficulty":3,"Vegetarian":true},{"RecipeID":3,"Name":"kofta curry","PrepTime":"10 mins ","Difficulty":3,"Vegetarian":true},{"RecipeID":11,"Name":"sambhar","PrepTime":"30 mins ","Difficulty":3,"Vegetarian":true},{"RecipeID":121,"Name":"fish pastry","PrepTime":"30 mins ","Difficulty":3,"Vegetarian":true},{"RecipeID":125,"Name":"fish moussaka","PrepTime":"30 mins ","Difficulty":3,"Vegetarian":true}]`

3. List a recipe based on recipe ID:   
           URL: `/recipes/{id}`
           Method Type: `GET` 
           Headers: None 
           Sample Response: `{"RecipeID":1,"Name":"roti","PrepTime":"10 mins ","Difficulty":2,"Vegetarian":true}`
          

4. Create a Recipe:   
           URL: `/recipes`
           Method Type: `POST` 
           Headers: Authorization: Bearer  
           Request Body: 
           `{"Name": "apple stew", "PrepTime": "60 mins", "Difficulty": 2, "Vegetarian": true}`
           Sample Response:`{"Result":"Successfully created"}`    

5. Modify a Recipe based on recipe ID:   
           URL: `/recipes/{id}`
           Method Type: `PUT` 
           Headers: `Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjIwMDAwLCJ1c2VybmFtZSI6ImFkbWluIn0` 
           Request Body: 
            `{"Name": "apple stew", "PrepTime": "60 mins", "Difficulty": 5, "Vegetarian": true}`
           Sample Response:  `{"result": "Recipe successfully modfied"}`          

6. Delete a Recipe based on recipe ID:   
           URL: `/recipes/{id}`
           Method Type: `DELETE` 
           Headers: `Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjIwMDAwLCJ1c2VybmFtZSI6ImFkbWluIn0` 
           Sample Response: `{"result":"Recipe successfully deleted"}`

7. Rate a Recipe based on recipe ID:   
           URL: `/recipes/{id}/rating`
           Method Type: `POST` 
           Headers: `Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjIwMDAwLCJ1c2VybmFtZSI6ImFkbWluIn0` 
           Request Body: `{"Rating": 5}`
           Sample Response: `{"result":"Recipe Ratings added"}`

8. Search a recipe based on keyword, the api will return recipes which contain the queried keyword:   
           URL: `/recipes/search/?search=<search_query>`
           Method Type: `GET` 
           Headers: None
           Sample Response: `[{"RecipeID":7,"Name":"cheesecake","PrepTime":"120 mins","Difficulty":4,"Vegetarian":true}]`            

##### Running the tests

```
ginkgo -r
```
##### Future Improvements/ Enhancements   

- [ ] Introduce middleware for authentication

- [ ] Remove hardcoded credentials 

- [ ] Better request validation 

- [ ] keyset pagination can be implemented, which is more efficient 

- [ ] Error handling and error messages can be made more robust. 

- [ ] Higher test case coverage, Negative cases  

- [ ] Content negotiation  

- [ ] Implementation of a cache to improve performance. 

- [ ] Gorm and database abstraction

 


