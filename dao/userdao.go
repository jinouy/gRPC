package dao

import (
	"gRPC_User/model"
	"gRPC_User/util"
	"github.com/pkg/errors"
)

func GetUserByName(u *model.User) (model.User, error) {

	var users model.User

	user := util.Db.Model(&users).Where("name = ?", u.Name)
	if users.Name != u.Name {
		return model.User{}, errors.Errorf("名字不存在")
	}
	if user.Error != nil {
		return model.User{}, user.Error
	}

	return users, nil
}

func GetUsersPage(p *model.Page) ([]model.User, error) {

	var users []model.User
	var total int64

	if err := util.Db.Model(&model.User{}).Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, err
	}

	if err := util.Db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, err
	}
	if int64(p.PageSize) > total {
		var data []byte
		data = []byte("查询超出数据长度")
		return nil, errors.Errorf(string(data))
	}
	return users, nil

}

func AddUser(u *model.User) (*model.User, error) {

	user := &model.User{Name: u.Name}

	users := util.Db.Create(&user)
	if users.Error != nil {
		return nil, users.Error
	}
	return user, nil
}

func DpdUser(u *model.Username) (model.User, error) {

	var user model.User

	users := util.Db.Model(&user).Where("name = ?", u.OldName).Update("name", u.NewName)
	if users.Error != nil {
		return model.User{}, users.Error
	}

	return user, nil
}

func DelUser(u *model.User) error {

	var users model.User

	user := util.Db.Where("name = ?", u.Name).Delete(&users)
	if user.Error != nil {
		return user.Error
	}
	return nil
}
