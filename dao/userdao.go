package dao

import (
	"gRPC_User/model"
	"gRPC_User/util"
)

func GetUserByName(u *model.User) (*model.User, error) {

	var user model.User

	DB := util.Db.Where("name = ?", u.Name).First(&user)
	if DB.Error != nil {
		return nil, DB.Error
	}

	return &user, nil
}

func GetUsersPage(p *model.Page) ([]*model.User, error) {

	var users []*model.User
	var total int64

	if err := util.Db.Model(&model.User{}).Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, err
	}

	if err := util.Db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, err
	}
	return users, nil

}

func AddUser(u *model.User) (*model.User, error) {

	user := &model.User{Name: u.Name}

	DB := util.Db.Create(&user)
	if DB.Error != nil {
		return nil, DB.Error
	}
	return user, nil
}

func DpdUser(u *model.Username) (*model.User, error) {

	var user model.User

	if err := util.Db.Model(&user).Where("name = ?", u.OldName).Update("name", u.NewName).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func DelUser(u *model.User) error {

	var users model.User

	if err := util.Db.Where("name = ?", u.Name).Delete(&users).Error; err != nil {
		return err
	}
	return nil
}
