package crud


import (
	"context"
	"charitable/model"
	"net/http"
	"encoding/json"
	"strconv"
	"io"
	"github.com/go-chi/chi"
		
	)

    type Food interface {
		GetFoods(ctx context.Context, limit int64) ([]*model.OrgFood, error)
		GetFoodsByOrgID(ctx context.Context, orgID int64) (*[]model.OrgFood, error)
		GetOrgFood(ctx context.Context, orgID int64, categoryID int64) (*[]model.OrgFood, error)
		GetFoodsByID(ctx context.Context, FoodID int64) (*model.OrgFood, error)
		DeleteFood(ctx context.Context, FoodID int64) (bool, error)
		DeleteFoodInfo(ctx context.Context, FoodID int64) (bool, error)

	
	}

	func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
		response, _ := json.Marshal(payload)
	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(response)
	}


func (repo *Repo) getFoodSizeByID(foodSizeID int) (*model.FoodSize ,error){
	
	
			rows, err := repo.DBConn.Query("SELECT food_size_id, name FROM public.food_size WHERE food_size_id = $1",foodSizeID)
			if err != nil {
				logger.Logger("ERROR", "gEetFoodSizeByID", err)
				panic(err)
			}
			defer rows.Close()

			food := new(model.FoodSize)
			for rows.Next() {	
			rows.Scan(&food.FoodSizeID, &food.Name)
			}

	return food , err
	}
		




func (repo *Repo)  getFoodInfoByFoodID(foodID int) ([]*model.FoodInfo, error){
	
	
	rows, err := repo.DBConn.Query("SELECT food_info_id, org_food_id ,cook_time, price ,food_size_id,packing_time FROM public.food_info WHERE org_food_id = $1",foodID)
	if err != nil {
		logger.Logger("ERROR", "gEetFoodSizeByID", err)
		panic(err)
	}
	defer rows.Close()

	var food_size_id int
	payload := make([]*model.FoodInfo, 0)

	
	for rows.Next() {
		food := new(model.FoodInfo)	
		food_size := new(model.FoodSize)
		err := rows.Scan(
			&food.FoodInfoID,
			&food.OrgFoodID,
			&food.CookTime,
			&food.Price,
			&food_size_id,
			&food.PackingTime,
		)
		if err == nil{
			food_size , err = repo.getFoodSizeByID(food_size_id)
			if err == nil{
				food.Size = model.FoodSize{FoodSizeID : food_size.FoodSizeID, Name : food_size.Name}
			}
		}
		if err!=nil{
			logger.Logger("ERROR", "getFoodInfoByFoodID", err)
				panic(err)
		}
		payload = append(payload, food)
	}
	
	return payload, err
}





func (repo *Repo) GetFood(ctx context.Context, org_id int64, food_id int64) (*model.OrgFood){

	rows, err := repo.DBConn.Query("SELECT org_food_id, org_id ,foodname, foods_category_id , description FROM public.org_foods WHERE  org_food_id= $1", food_id)
	
	if err != nil {
		logger.Logger("ERROR", "GetOrgFood", err)
				panic(err)
	}

	defer rows.Close()
	
	//foods := make([]*model.OrgFood, 0)
	food := new(model.OrgFood)

	if rows.Next() {	
	
	
	foodInfo := make([]*model.FoodInfo,0)
	
	err := rows.Scan(&food.OrgFoodID, &food.OrgID,&food.FoodName,&food.FoodCategory,&food.FoodDescription)

	if err==nil {
		foodInfo , err = repo.getFoodInfoByFoodID(food.OrgFoodID)

		if err == nil{
			food.FoodsInfo = foodInfo
		}
		
	}
	
	photos,err:= repo.GetPhotos(food.OrgFoodID,"F")
	if err!=nil{
		logger.Logger("ERROR","GetOrgPhoto",err)
	}
	for _,i:=range photos{
		food.ImageUrl = append(food.ImageUrl,i)
	}
	
	
	if err != nil {
		logger.Logger("ERROR", "GetOrgFood", err)
				panic(err)
	}
	
    }
	return food;
}






func (repo *Repo) CreateFood(food model.OrgFood )(model.OrgFood,error){

		rows ,err:=repo.DBConn.Query("INSERT INTO public.org_foods(org_id ,foodname,foods_category_id,description) VALUES($1,$2,$3,$4) RETURNING org_food_id",food.OrgID,food.FoodName,food.FoodCategory,food.FoodDescription)
		
		if err!=nil{
			logger.Logger("ERROR", "CreatFood", err)
			panic(err)
		}
		rows.Next()
		rows.Scan(&food.OrgFoodID)
		println("orgfooid: ",food.OrgFoodID)
		defer rows.Close()
			
		print("len: ",len(food.FoodsInfo))
		for _, f := range food.FoodsInfo {
			f.OrgFoodID = food.OrgFoodID
			rows,err :=repo.DBConn.Query("INSERT INTO public.food_info(org_food_id ,food_size_id,cook_time, price,packing_time) VALUES($1,$2,$3,$4,$5)RETURNING food_info_id",f.OrgFoodID, f.Size.FoodSizeID , f.CookTime, f.Price,f.PackingTime) 
			t,err:=repo.getFoodSizeByID(f.Size.FoodSizeID)
			if err !=nil{
				return food,err
			}
			f.Size.Name= t.Name
			rows.Next()
			rows.Scan(&f.FoodInfoID)
			println("foodinfoid: ", f.FoodInfoID)
		}	
	
		defer rows.Close()
		food.ImageUrl,err= repo.GetPhotos(food.OrgFoodID,"F")
		if err!=nil{
			panic(err)
		}
	return food,nil
}



//Delete food
func (repo *Repo) DeleteFood(orgfoodID int64)(bool,error){		
		rows ,err:=repo.DBConn.Query("DELETE FROM public.food_info WHERE org_food_id = $1",orgfoodID)
		if err!=nil{
			panic(err)
		return false,err
	}
		
		rows ,err=repo.DBConn.Query("DELETE FROM public.org_foods WHERE org_food_id = $1",orgfoodID)
	//rows ,errr:=db.Query("insert into public.org_foods(org_id ,foodname,foods_category_id) values($1,$2,$3) RETURNING org_food_id",food.Org_id,food.FoodName,food.FoodCategory).Scan(&food.OrgFood_id)
	
	rows.Next()
return true ,nil
}
func (repo *Repo) DeleteFoodInfo(foodinfoID int64)(bool,error){		
	rows ,err:=repo.DBConn.Query("DELETE FROM public.food_info WHERE food_info_id = $1",foodinfoID)
	if err!=nil{
		panic(err)
	return false,err
}
rows.Next()
return true ,nil
}


func (repo *Repo) UpdateFood(food model.OrgFood, food_id string,org_id int64 )(model.OrgFood,error){

	rows ,err:=repo.DBConn.Query("UPDATE public.org_foods SET foodname=$1,foods_category_id =$2, description= $3 where org_food_id=$4 and org_id= $5",food.FoodName,food.FoodCategory,food.FoodDescription,food_id,org_id)
	//rows ,errr:=db.Query("insert into public.org_foods(org_id ,foodname,foods_category_id) values($1,$2,$3) RETURNING org_food_id",food.Org_id,food.FoodName,food.FoodCategory).Scan(&food.OrgFood_id)
	if err!=nil{
		println("hes error")
		panic(err)
	}
	rows.Next()
	rows.Scan(&food.OrgFoodID)
	println("orgfooid: ",food.OrgFoodID)
	defer rows.Close()
		
	print("len: ",len(food.FoodsInfo))
	for _, f := range food.FoodsInfo {
		f.OrgFoodID = food.OrgFoodID
		print(f.Size.FoodSizeID ," ", f.CookTime," ", f.Price, " ",f.PackingTime," ",f.FoodInfoID)
		rows,err :=repo.DBConn.Query("UPDATE public.food_info SET food_size_id=$1, cook_time= $2, price=$3,packing_time=$4 WHERE food_info_id=$5" ,f.Size.FoodSizeID , f.CookTime, f.Price,f.PackingTime,f.FoodInfoID) 
		t,err:=repo.getFoodSizeByID(f.Size.FoodSizeID)
		if err !=nil{
			return food,err
		}
		f.Size.Name= t.Name
		rows.Next()
		rows.Scan(&f.FoodInfoID)
		println("foodinfoid: ", f.FoodInfoID)
	}	

	
	
	

return food,nil
}


		
