package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// elastic "gopkg.in/olivere/elastic"
)

var DB *gorm.DB
var err error
var userCache UserCache = NewRedisCache("localhost:6379", 0, 10)

// esclient, err := GetESClient()
// if err != nil {
// 	fmt.Println("Error initializing : ", err)
// 	panic("Client fail ")
// }

const DNS = "host=localhost user=girishjain password=12345 dbname=godb port=5432 sslmode=disable TimeZone=Asia/Shanghai"

type User struct {
	gorm.Model
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}

func InitialMigration() {
	DB, err = gorm.Open(postgres.Open(DNS), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	}
	DB.AutoMigrate(&User{})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User
	DB.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userID := params["id"]
	var user *User = userCache.Get(userID)
	if user == nil {
		DB.First(&user, userID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("No posts found")
			return
		}
		userCache.Set(userID, *user)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	DB.Create(&user)
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user User
	DB.First(&user, params["id"])
	json.NewDecoder(r.Body).Decode(&user)
	DB.Save(&user)
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user User
	DB.Delete(&user, params["id"])
	json.NewEncoder(w).Encode("The USer is Deleted Successfully!")
}

// func FindByID(id string) (*User, error) {

// 	row := DB.Query("select * from users where id = ?", id)

// 	var user User
// 	if row != nil {
// 		var id int64
// 		var title string
// 		var text string
// 		err := row.Scan(&id, &title, &text)
// 		if err != nil {
// 			return nil, err
// 		} else {
// 			user = User{
// 				FirstName: firstname,
// 				LastName: lastname,
// 				Email :  email,
// 			}
// 		}
// 	}

// 	return &user, nil
// }
