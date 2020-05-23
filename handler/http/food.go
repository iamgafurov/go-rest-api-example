package handler

import (
	//"google.golang.org/api/drive/v3"   
	"io"
	driver "charitable/driver"
	//logger "charitable/logger"
	//service "charitable/service"
	//service "charitable/service"
	"charitable/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	//"io/ioutil"
	"github.com/go-chi/chi"
)

func (p *OrgHandler) Auth(w http.ResponseWriter, r *http.Request) {

	print(r.Method)

	if r.Method != "POST" {
	
		http.Error(w, http.StatusText(405), 405)
		return
	}

	err := r.ParseForm()

	if err != nil {
		panic(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	
	authInfo, err := p.Repo.Autorization(r.Context(),email,password)

	println("auth info :" , authInfo.Token)
	  if err == nil{
		RespondWithJSON(w, http.StatusOK, authInfo)
	  }else{
		  RespondWithMessage(w, http.StatusNotFound,"error")
	  }

	
}

//Create Foods
	func (p *OrgHandler) CreateFoods(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {

		
			http.Error(w, http.StatusText(405), 405)
	
			return
		}
	food :=model.OrgFood{}
	
	org_id ,err:= p.Repo.DecodeToken(r.Header.Get("token"))
	if err !=nil{
		panic(err)	
	}
json.NewDecoder(r.Body).Decode(&food)
println("food_name: ",food.FoodName)
	food.OrgID = org_id
	food,err =p.Repo.CreateFood(food)
	if err!= nil{
		logger.Logger("ERROR","CreateFood food added orgfoodID: " + strconv.Itoa(food.OrgFoodID),err)
		panic(err)
	} 
	
	dir,err:= service.GetFoldersByName(strconv.FormatInt(food.OrgID,10))

	if err !=nil {
	 logger.Logger("ERROR","Creatfood",err)
	 println(err)
	}
	//dir,err:= service.GetFoldersByName(strconv.FormatInt(food.OrgID,10))
	//dir,_=service.GetFoldersByNameAndParentId("foods",dir.Id)

	//dir,_= service.CreateDir(strconv.Itoa(food.OrgFoodID),dir.Id)
	
		//_,err=p.Repo.DBConn.Query("UPDATE public.org_foods SET images_url = $1 WHERE org_food_id= $2",dir.Id, food.OrgFoodID)
		
		//if err!=nil{
		//	println(err)
		//}

	
	RespondWithJSON(w , 200,food)
	
	
}

//Update Foods
func (p *OrgHandler) UpdateFoods(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {

	
		http.Error(w, http.StatusText(405), 405)

		return
	}
//org_id := chi.URLParam(r, "org_id")
food :=model.OrgFood{}
food_id := chi.URLParam(r,"food_id")
println(food_id)
org_id ,err:= p.Repo.DecodeToken(r.Header.Get("token"))
if err !=nil{
	panic(err)	
}
json.NewDecoder(r.Body).Decode(&food)
println("org_id: ",org_id)
food.OrgID = org_id
food,err =p.Repo.UpdateFood(food , food_id,org_id)
if err!= nil{
	panic(err)
} 
RespondWithJSON(w , 200,food)


}
	


//GetFoods
func (p *OrgHandler) GetFoods(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
//org_id := chi.URLParam(r, "org_id")
	println(r.Header.Get("token"))
	var org_id,category_id int64
	org_id ,err:= p.Repo.DecodeToken(r.Header.Get("token"))
	if err !=nil{
		println("ERRR")
		panic(err)	
}


	println("org_id: ",org_id)
	fmt.Sscan(chi.URLParam(r, "category_id"), &category_id)
//category_id := chi.URLParam(r, "category_id")


	foods:= p.Repo.GetOrgFood(r.Context(), org_id, category_id)

	RespondWithJSON(w, http.StatusCreated, foods)


}

//GetFoodsByID
func (p *OrgHandler) GetFoodByID(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	//food_id := chi.URLParam(r, "food_id")

	var org_id, food_id int64

	org_id ,err:= p.Repo.DecodeToken(r.Header.Get("token"))
	if err !=nil{
		println("ERRR")
		panic(err)	
}




	println("org_id: ",org_id)
	
food_id,err= strconv.ParseInt(chi.URLParam(r, "food_id"),10,64)
//category_id := chi.URLParam(r, "category_id")
if err!=nil{
	panic(err)
}
println("org_food_id:",food_id)
	food := p.Repo.GetFood(r.Context(), org_id, food_id)

	RespondWithJSON(w, http.StatusCreated, food)


}


//Delete Food
func (p *OrgHandler) DeleteFood(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
//org_id := chi.URLParam(r, "org_id")
println(r.Header.Get("token"))
var orgfood_id int64
_ ,err:= p.Repo.DecodeToken(r.Header.Get("token"))
if err !=nil{
	panic(err)	
}

keys,ok := r.URL.Query()["orgfood_id"]
if !ok{
	panic(ok)
}
fmt.Sscan(keys[0],&orgfood_id)
println("orgfood_id: ",orgfood_id)

h,err := p.Repo.DeleteFood(orgfood_id)
if h!=true{
	RespondWithMessage(w,403,"ERROR")
}else{
	RespondWithMessage(w,200,"Success")
}



}


//Delete FoodInfo
func (p *OrgHandler) DeleteFoodInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
//org_id := chi.URLParam(r, "org_id")
	println(r.Header.Get("token"))
	var orgfood_id int64
	_ ,err:= p.Repo.DecodeToken(r.Header.Get("token"))
	if err !=nil{
		panic(err)	
	}else{

		keys,ok := r.URL.Query()["foodinfo_id"]
		if !ok{
			panic(ok)
		}
		fmt.Sscan(keys[0],&orgfood_id)
		println("foodinfo_id: ",orgfood_id)

		h,_ := p.Repo.DeleteFoodInfo(orgfood_id)
		if h!=true{
			RespondWithMessage(w,403,"ERROR")
		}else{
			RespondWithMessage(w,200,"Success")
		}
	}
}

func (p *OrgHandler) GooglePost(w http.ResponseWriter, r *http.Request) {


  println("url = ",r.URL,"Method =",r.Method)

	RespondWithMessage(w,200,"Success")
}