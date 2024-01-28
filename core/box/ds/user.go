package ds

import (
	"github.com/ipfs/kubo/core/box/model"
	"github.com/jinzhu/gorm"
)

func (s *DbStore) UserCount() (count int, _ error) {
	ret := s.db.Model(&model.User{}).Count(&count)
	return count, ret.Error
}

// func (s *DbStore) GetUserById(id int) (*model.User, error) {
// 	var user model.User
// 	ret := s.db.Model(&model.User{}).Where("id=?", id).First(&user)
// 	if ret.Error != nil {
// 		if ret.Error != gorm.ErrRecordNotFound {
// 			log.Errorf("%v", ret.Error.Error())
// 		}
// 	}
// 	return &user, ret.Error
// }

// func (s *DbStore) GetAllUsers() (list []*model.User, _ error) {

// 	ret := s.db.Model(&model.User{}).Order("id").Scan(&list)
// 	if ret.Error != nil {
// 		log.Errorf("%v", ret.Error.Error())
// 	}
// 	return list, ret.Error
// }

// func (s *DbStore) UpdateName(user *model.User) error {
// 	mm := map[string]interface{}{
// 		"name": user.Name,
// 	}
// 	return s.updateUser(user.Id, mm)
// }

func (s *DbStore) InserUserInfo(user model.User) error {
	result := s.db.Create(&user)
	if result.Error != nil {
		result = s.db.Save(&user)
		return result.Error
	}
	return nil
}

func (s *DbStore) GetUserInfo(Id string) (*model.User, error) {
	var user model.User
	ret := s.db.Model(&model.User{}).Where("id=?", Id).First(&user)
	if ret.Error != nil {
		if ret.Error != gorm.ErrRecordNotFound {
			log.Errorf("%v", ret.Error.Error())
		}
	}
	return &user, ret.Error
}
