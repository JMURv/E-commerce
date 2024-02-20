package models

import "time"

type Room struct {
	ID        uint      `gorm:"primaryKey"`
	Members   []string  `json:"members" gorm:"many2many:room_members"`
	ItemID    *uint     `json:"itemID"`
	Item      string    `json:"item" gorm:"foreignKey:ItemID"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"createdAt"`
}

func CreateRoom(user1ID, user2ID uint) (*Room, error) {
	//var user1, user2 User
	//
	//if err := db.First(&user1, user1ID).Error; err != nil {
	//	return nil, err
	//}
	//if err := db.First(&user2, user2ID).Error; err != nil {
	//	return nil, err
	//}
	//
	//newRoom := Room{
	//	Members:   []User{user1, user2},
	//	CreatedAt: time.Now(),
	//}
	//
	//if err := db.Create(&newRoom).Error; err != nil {
	//	return nil, err
	//}

	return &Room{}, nil
}

func GetUserRoomWithMessages(requestUserID uint) ([]Room, error) {
	var userRooms []Room

	if err := db.Preload("Members").Preload("Messages").Preload("Item").
		Joins("JOIN room_members ON room_members.room_id = rooms.id").
		Where("room_members.user_id = ?", requestUserID).
		Find(&userRooms).
		Error; err != nil {
		return nil, err
	}

	return userRooms, nil
}
