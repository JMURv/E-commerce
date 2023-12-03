package main

import (
	"e-commerce/pkg/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var db *gorm.DB

func ListCreateItems(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var items []models.Item
		result := db.Find(&items)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
		return
	} else if r.Method == "POST" {
		var newItem models.Item

		err := json.NewDecoder(r.Body).Decode(&newItem)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		result := db.Create(&newItem)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newItem)
	}
}

func GetItemByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	itemID := vars["id"]
	var item models.Item

	result := db.First(&item, itemID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func ListCreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		return
		//var categories []models.Category
		//result := db.Find(&categories)
	}
}

func init() {
	var err error
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	log.Println(".env file has been loaded")

	dsn := os.Getenv("DSN")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.Item{}, &models.Category{}, &models.Tag{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to the database")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/items", ListCreateItems).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/items/{id}", GetItemByID).Methods(http.MethodGet)

	log.Println("Server is running on port 8000")
	http.ListenAndServe(":8000", r)
}
