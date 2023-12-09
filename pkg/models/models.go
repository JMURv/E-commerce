package models

import (
	"e-commerce/pkg/config"
	"errors"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

type Item struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CategoryID  *uint     `json:"categoryID"`
	Category    *Category `json:"category" gorm:"foreignKey:CategoryID"`
	Sellers     []Seller  `json:"sellers" gorm:"many2many:seller_items;"`
	Tags        []Tag     `json:"tags" gorm:"many2many:item_tags;"`
}

type Category struct {
	gorm.Model
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	ParentCategoryID *uint     `json:"parentCategoryID"`
	ParentCategory   *Category `json:"parentCategory" gorm:"foreignKey:ParentCategoryID"`
}

type User struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Favorite struct {
	gorm.Model
	Author         User `json:"author" gorm:"foreignKey:AuthorID"`
	AuthorID       uint `json:"authorID"`
	FavoriteItem   Item `json:"favoriteItem" gorm:"foreignKey:FavoriteItemID"`
	FavoriteItemID uint `json:"favoriteItemID"`
}

type Seller struct {
	gorm.Model
	Username  string `json:"username"`
	Email     string `json:"email"`
	SoldItems []Item `json:"soldItems" gorm:"many2many:seller_items;"`
}

type Cart struct {
	gorm.Model
	UserID uint       `json:"userID"`
	Items  []CartItem `json:"items" gorm:"many2many:cart_items;"`
}

type CartItem struct {
	CartID   uint `gorm:"primaryKey"`
	ItemID   uint `gorm:"primaryKey"`
	Quantity int  `json:"quantity"`
}

type Order struct {
	gorm.Model
	UserID   uint        `json:"userID"`
	Items    []OrderItem `json:"items" gorm:"many2many:order_items;"`
	Status   string      `json:"status"`
	Payment  string      `json:"payment"`
	Shipping string      `json:"shipping"`
}

type OrderItem struct {
	OrderID  uint `gorm:"primaryKey"`
	ItemID   uint `gorm:"primaryKey"`
	Quantity int  `json:"quantity"`
}

type Review struct {
	gorm.Model
	Author         User   `json:"author" gorm:"foreignKey:AuthorID"`
	AuthorID       uint   `json:"authorID"`
	ReviewedItem   Item   `json:"reviewedItem" gorm:"foreignKey:ReviewedItemID"`
	ReviewedItemID uint   `json:"reviewedItemID"`
	Advantages     string `json:"advantages"`
	Disadvantages  string `json:"disadvantages"`
	ReviewText     string `json:"reviewText"`
	Rating         int    `json:"rating"`
}

type Tag struct {
	gorm.Model
	Name string `json:"name"`
}

func init() {
	var err error
	config.Connect()
	db = config.GetDB()

	err = db.AutoMigrate(&Item{}, &Category{}, &Tag{}, &User{}, &Seller{}, &Cart{}, &CartItem{}, &Order{}, &OrderItem{}, &Favorite{})
	if err != nil {
		log.Fatal(err)
	}
}

func (i *Item) CreateItem() (*Item, error) {

	if i.Name == "" {
		return i, errors.New("name is required")
	}
	if i.Description == "" {
		return i, errors.New("description is required")
	}
	if i.Price == 0 {
		return i, errors.New("price is required")
	}

	for idx := range i.Tags {
		existingTag := &Tag{}
		if err := db.Where("name = ?", i.Tags[idx].Name).First(existingTag).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&i.Tags[idx]).Error; err != nil {
					return nil, err
				}
				if err := db.Where("name = ?", i.Tags[idx].Name).First(existingTag).Error; err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		i.Tags[idx] = *existingTag
	}

	result := db.Create(&i)
	if result.Error != nil {
		return nil, result.Error
	}
	return i, nil
}

func (i *Item) UpdateItem(newData *Item) (*Item, error) {
	if newData.Name != "" {
		i.Name = newData.Name
	}

	if newData.Description != "" {
		i.Description = newData.Description
	}

	if newData.Price != 0 {
		i.Price = newData.Price
	}

	if newData.CategoryID != i.CategoryID {
		i.CategoryID = newData.CategoryID
	}
	db.Save(&i)
	return i, nil
}

func UpdateUser(userId string, newData *User) (*User, error) {
	user := GetUserByID(userId)
	if newData.Username != "" {
		user.Username = newData.Username
	}

	if newData.Email != "" {
		user.Email = newData.Email
	}
	db.Save(&user)
	return user, nil
}

func (c *Category) CreateCategory() (*Category, error) {
	if c.Name == "" {
		return c, errors.New("name is required")
	}
	if c.Description == "" {
		return c, errors.New("description is required")
	}
	result := db.Create(&c)
	if result.Error != nil {
		return nil, result.Error
	}
	return c, nil
}

func (u *User) CreateUser() (*User, error) {
	if u.Username == "" {
		return u, errors.New("username is required")
	}
	if u.Email == "" {
		return u, errors.New("email is required")
	}
	result := db.Create(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return u, nil
}

func GetAllItems() []Item {
	var Items []Item
	db.Preload("Tags").Find(&Items)
	return Items
}

func GetAllCategories() []Category {
	var Categories []Category
	db.Preload("ParentCategory").Find(&Categories)
	return Categories
}

func GetAllUsers() []User {
	var Users []User
	db.Find(&Users)
	return Users
}

func GetItemByID(id string) *Item {
	var getItem Item
	db.Preload("Tags").Where("ID=?", id).First(&getItem)
	return &getItem
}

func GetUserByID(id string) *User {
	var user User
	db.Where("ID=?", id).First(&user)
	return &user
}

func DeleteItem(id string) Item {
	var item Item
	db.Delete(&item, id)
	return item
}

func DeleteUser(id string) User {
	var user User
	db.Delete(&user, id)
	return user
}
