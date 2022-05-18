package service

import (
	"context"
	"gRPC_User/dao"
	"gRPC_User/model"
	"gRPC_User/proto"
	"log"
	"net/http"
)

//var UserService = &userService{}

type UserService struct {
	proto.UnimplementedUserServiceServer
}

func (p *UserService) GetUserByName(ctx context.Context, req *proto.UserGetReq) (*proto.UserGetResp, error) {

	u := &model.User{Name: req.Name}

	user, err := dao.GetUserByName(u)
	if err != nil {
		return &proto.UserGetResp{Code: http.StatusInternalServerError, Msg: err.Error()}, nil
	}

	users := &proto.User{
		Id:        int64(user.ID),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		Name:      user.Name,
	}

	return &proto.UserGetResp{User: users, Code: http.StatusOK, Msg: "获取数据成功"}, nil

}

func (p *UserService) GetUsers(ctx context.Context, req *proto.UsersGetReq) (*proto.UsersGetResp, error) {

	if len(string(req.Page)) == 0 || len(string(req.Limit)) == 0 {
		return &proto.UsersGetResp{Code: http.StatusBadRequest, Msg: "参数为空"}, nil
		log.Fatal("参数为空")
	}

	var re proto.UsersGetResp

	page := &model.Page{
		PageNum:  int(req.Page),
		PageSize: int(req.Limit),
	}

	users, err := dao.GetUsersPage(page)
	if err != nil {
		return &proto.UsersGetResp{Code: http.StatusInternalServerError, Msg: err.Error()}, nil
	}
	if len(users) == 0 {
		return &proto.UsersGetResp{Code: http.StatusInternalServerError, Msg: "查询页面过大"}, nil
	}

	for _, v := range users {
		user := &proto.User{
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

func (p *UserService) AddUser(ctx context.Context, req *proto.UserPostReq) (*proto.UserPostResp, error) {

	if len(req.Name) < 2 || len(req.Name) > 8 {
		return &proto.UserPostResp{Code: http.StatusBadRequest, Msg: "参数长度不正确"}, nil
		log.Fatal("参数长度不正确")
	}
	u := &model.User{Name: req.Name}

	Userknown, err := dao.GetUserByName(u)
	if Userknown.Name == req.Name {
		return &proto.UserPostResp{Code: http.StatusBadRequest, Msg: "名字已存在"}, nil
	}
	if err != nil {
		if err.Error() == "名字不存在" {
			users, err := dao.AddUser(u) //进行数据库的增加操作

			if err != nil { //数据库怠机，返回错误
				return &proto.UserPostResp{Code: http.StatusInternalServerError, Msg: err.Error()}, err
			}

			user := &proto.User{
				Id:        int64(users.ID),
				CreatedAt: users.CreatedAt.String(),
				UpdatedAt: users.UpdatedAt.String(),
				Name:      users.Name,
			}
			return &proto.UserPostResp{User: user, Code: http.StatusCreated, Msg: "数据添加成功"}, nil

		} else {
			return &proto.UserPostResp{Code: http.StatusInternalServerError, Msg: err.Error()}, err
		}
	}
	return nil, nil

}

func (p *UserService) UpdUserName(ctx context.Context, req *proto.UserPutReq) (*proto.UserPutResp, error) {

	if len(req.NewName) < 2 || len(req.NewName) > 10 || len(req.OldName) == 0 || len(req.NewName) == 0 {
		return &proto.UserPutResp{Code: http.StatusBadRequest, Msg: "参数错误"}, nil
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
	if errUnknown != (model.User{}) {
		return &proto.UserPutResp{Code: http.StatusBadRequest, Msg: "名字已存在"}, nil
	}
	if err != nil {
		return &proto.UserPutResp{Code: http.StatusInternalServerError, Msg: "数据库怠机"}, err
	}
	_, err = dao.GetUserByName(userOldname)
	if err != nil { //添加相同名字的限制条件，如果相同就返回错误
		return &proto.UserPutResp{Code: http.StatusBadRequest, Msg: err.Error()}, nil
	}
	if err != nil {
		return &proto.UserPutResp{Code: http.StatusInternalServerError, Msg: "数据库怠机"}, err
	}
	users, err := dao.DpdUser(username)
	if err != nil {
		return &proto.UserPutResp{Code: http.StatusInternalServerError, Msg: "数据库怠机"}, err
	}

	user := &proto.User{
		Id:        int64(users.ID),
		CreatedAt: users.CreatedAt.String(),
		UpdatedAt: users.UpdatedAt.String(),
		Name:      users.Name,
	}
	return &proto.UserPutResp{User: user, Code: http.StatusCreated, Msg: "插入成功"}, nil
}

func (p *UserService) DelUser(context context.Context, req *proto.UserDelReq) (*proto.UserDelResp, error) {

	u := &model.User{
		Name: req.Name,
	}
	users, err := dao.GetUserByName(u)
	if err != nil {
		return &proto.UserDelResp{Code: http.StatusBadRequest, Msg: err.Error()}, nil
	}
	err = dao.DelUser(u)
	if err != nil {
		return &proto.UserDelResp{Code: http.StatusInternalServerError, Msg: "数据库怠机"}, err
	}

	user := &proto.User{
		Id:        int64(users.ID),
		CreatedAt: users.CreatedAt.String(),
		UpdatedAt: users.UpdatedAt.String(),
		Name:      users.Name,
	}

	return &proto.UserDelResp{User: user, Code: http.StatusCreated, Msg: "删除成功"}, nil

}
