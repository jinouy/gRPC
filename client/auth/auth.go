package auth

import (
	"fmt"
	pb "gRPC_User/proto/chat"
)

var User = InputName()

func InputName() string {
	var baseMsg pb.HiRequest
	fmt.Println("请输入用户昵称：")
	_, _ = fmt.Scanln(&baseMsg.Name)
	return baseMsg.Name
}

//func Clientinerceptor(ctx context.Context, method string, req, reply interface{},
//	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
//
//	start := time.Now()
//	err := invoker(ctx, method, req, reply, cc, opts...)
//	log.Printf("method == %s ; req == %v ; rep == %v ; duration == %s ; error == %v\n", method, req, reply, time.Since(start), err)
//	return err
//}
