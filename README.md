# RESTful API Example with golang
This is simple example restful api server 

## Install and Run
```shell
$ go get github.com/iamgafurov/go-rest-api-example

$ cd $GOPATH/src/github.com/iamgafurov/go-rest-api-example
$ go build
$ ./go-rest-api-example
```

## API Endpoint
- http://localhost:2222/api/v1/foods/{category}
    - `GET`: get list of foods by category
    - `POST`: create food
- http://localhost:2222/api/v1/foods/food/{id}
    - `GET`: get food by id
    - `PUT`: update food
    - `DELETE`: remove food

   http://localhost:2222/api/v1/foods/foodinfo/
    
    - `DELETE`: remove foodinfo
    
## Data Structure
```json
 {
	"FoodInfoID": 	int ,   		
	"OrgFoodID" : 	int	,			
	"Size"			:   FoodSize ,		
	"CookTime" 	:	  string ,			
	"PackingTime":	string	,		
	"Price"     : 			float32		
}

FoodSize struct {
	FoodSizeID 		int			
	Name 			string			
}

```
