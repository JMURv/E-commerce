package model

type Tag struct {
	Name string `json:"name" gorm:"primaryKey"`
}

//
//func GetAllTags() ([]Tag, error) {
//	var tags []Tag
//	if err := db.Find(&tags).Error; err != nil {
//		return nil, err
//	}
//	return tags, nil
//}
//
//func (t *Tag) CreateTag() (*Tag, error) {
//	if t.Name == "" {
//		return nil, errors.New("name is required")
//	}
//
//	if err := db.Create(&t).Error; err != nil {
//		return nil, err
//	}
//	return t, nil
//}
//
//func DeleteTag(tagID uint) error {
//	var tag Tag
//	if err := db.Unscoped().Delete(&tag, tagID).Error; err != nil {
//		return err
//	}
//	return nil
//}
