package controller

import (
	"fmt"
	pb "gRPC_User/proto/chat"
)

func InputName() string {

	var baseMsg pb.HiRequest
	fmt.Println("请输入用户昵称：")
	_, _ = fmt.Scanln(&baseMsg.Name)
	fmt.Println("用户昵称为：", baseMsg.Name)

	return baseMsg.Name
}
