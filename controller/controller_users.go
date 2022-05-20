package controller

import (
	"context"
	"gRPC_User/dao"
	"gRPC_User/model"
	"gRPC_User/proto/user"
	"gorm.io/gorm"
	"log"
	"net/http"
)

//var UserService = &userService{}

type UserService struct {
	user.UnimplementedUserServiceServer
}

func (p *UserService) GetUserByName(ctx context.Context, req *user.UserGetReq) (*user.UserGetResp, error) {

	u := &model.User{Name: req.Name}

	userpro, err := dao.GetUserByName(u)
	if err != nil {
		return &user.UserGetResp{Code: http.StatusInternalServerError, Msg: err.Error()}, nil
	}

	users := &user.User{
		Id:        int64(userpro.ID),
		CreatedAt: userpro.CreatedAt.String(),
		UpdatedAt: userpro.UpdatedAt.String(),
		Name:      userpro.Name,
	}

	return &user.UserGetResp{User: users, Code: http.StatusOK, Msg: "获取数据成功"}, nil

}

func (p *UserService) UserList(ctx context.Context, req *user.ListGetReq) (*user.ListGetResp, error) {

	if len(string(req.Page)) == 0 || len(string(req.Limit)) == 0 {
		return &user.ListGetResp{Code: http.StatusBadRequest, Msg: "参数为空"}, nil
		log.Fatal("参数为空")
	}

	var re user.ListGetResp

	page := &model.Page{
		PageNum:  int(req.Page),
		PageSize: int(req.Limit),
	}

	users, err := dao.GetUsersPage(page)
	if err != nil {
		return &user.ListGetResp{Code: http.StatusInternalServerError, Msg: err.Error()}, nil
	}
	if len(users) == 0 {
		return &user.ListGetResp{Code: http.StatusInternalServerError, Msg: "查询页面过大"}, nil
	}

	for _, v := range users {
		user := &user.User{
			Id:        int64(v.ID),
			CreatedAt: v.CreatedAt.String(),
			UpdatedAt: v.UpdatedAt.String(),
			Name:      v.Name,
		}
		re.User = append(re.User, user)
	}
	re.Code = http.StatusOK
	re.Msg = "获取数据成功"

	return &re, nil

}

func (p *UserService) AddUser(ctx context.Context, req *user.UserAddReq) (*user.UserAddResp, error) {

	if len(req.Name) < 2 || len(req.Name) > 8 {
		return &user.UserAddResp{Code: http.StatusBadRequest, Msg: "参数长度不正确"}, nil
		log.Fatal("参数长度不正确")
	}
	u := &model.User{Name: req.Name}

	Known, err := dao.GetUserByName(u)
	if err != nil && err != gorm.ErrRecordNotFound {
		return &user.UserAddResp{Code: http.StatusInternalServerError, Msg: err.Error()}, err
	}
	if Known != nil {

		return &user.UserAddResp{Code: http.StatusInternalServerError, Msg: "名字已经存在"}, err
	}

	users, err := dao.AddUser(u) //进行数据库的增加操作
	if err != nil {              //数据库怠机，返回错误
		return &user.UserAddResp{Code: http.StatusInternalServerError, Msg: err.Error()}, err
	}
	userpro := &user.User{
		Id:        int64(users.ID),
		CreatedAt: users.CreatedAt.String(),
		UpdatedAt: users.UpdatedAt.String(),
		Name:      users.Name,
	}
	return &user.UserAddResp{User: userpro, Code: http.StatusCreated, Msg: "数据添加成功"}, nil
}

func (p *UserService) UpdUserName(ctx context.Context, req *user.UserUpdReq) (*user.UserUpdResp, error) {

	if len(req.NewName) < 2 || len(req.NewName) > 10 || len(req.OldName) == 0 || len(req.NewName) == 0 {
		return &user.UserUpdResp{Code: http.StatusBadRequest, Msg: "参数错误"}, nil
		log.Fatal("参数错误")
	}

	username := &model.Username{ //将结构体u中的值传入定义的username中
		OldName: req.OldName,
		NewName: req.NewName,
	}
	userNewname := &model.User{
		Name: req.NewName,
	}
	userOldname := &model.User{
		Name: req.OldName,
	}

	errUnknown, err := dao.GetUserByName(userNewname)

	if errUnknown != nil {

		return &user.UserUpdResp{Code: http.StatusInternalServerError, Msg: "名字已存在"}, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {

		return &user.UserUpdResp{Code: http.StatusInternalServerError, Msg: err.Error()}, err
	}
	errknown, err := dao.GetUserByName(userOldname)
	if errknown != nil {

		users, err := dao.DpdUser(username)
		if err != nil {
			return &user.UserUpdResp{Code: http.StatusInternalServerError, Msg: "数据库怠机"}, err
		}

		userpro := &user.User{
			Id:        int64(users.ID),
			CreatedAt: users.CreatedAt.String(),
			UpdatedAt: users.UpdatedAt.String(),
			Name:      users.Name,
		}
		return &user.UserUpdResp{User: userpro, Code: http.StatusCreated, Msg: "插入成功"}, nil
	} else {
		return &user.UserUpdResp{Code: http.StatusInternalServerError, Msg: "名字不存在"}, nil
	}

}

func (p *UserService) DelUser(context context.Context, req *user.UserDelReq) (*user.UserDelResp, error) {

	u := &model.User{
		Name: req.Name,
	}
	users, err := dao.GetUserByName(u)
	if err != nil {
		return &user.UserDelResp{Code: http.StatusInternalServerError, Msg: err.Error()}, nil
	}
	err = dao.DelUser(u)
	if err != nil {
		return &user.UserDelResp{Code: http.StatusInternalServerError, Msg: "数据库怠机"}, err
	}

	userpro := &user.User{
		Id:        int64(users.ID),
		CreatedAt: users.CreatedAt.String(),
		UpdatedAt: users.UpdatedAt.String(),
		Name:      users.Name,
	}

	return &user.UserDelResp{User: userpro, Code: http.StatusCreated, Msg: "删除成功"}, nil

}
