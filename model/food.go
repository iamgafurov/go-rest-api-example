package model


type FoodSize struct {
	FoodSizeID 		int				`json:foodsize_id`
	Name 			string			`json:name`
}

type FoodInfo struct {
	FoodInfoID  	int     		`json: foodInfoid`
	OrgFoodID 		int				`json  orgfoodid`
	Size			FoodSize		`json: size`
	CookTime 		string			`json: cooktime`
	PackingTime 	string			`json: packingtime`
	Price 			float32			`json: price`
}
