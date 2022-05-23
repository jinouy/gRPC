package main

//func main() {
//
//	rpcServer := grpc.NewServer(grpc.Creds(comm.GetCertService()), grpc.UnaryInterceptor(comm.GetToke()))
//
//	src := &controller.UserService{}
//	user.RegisterUserServiceServer(rpcServer, src)
//
//	listener, err := net.Listen("tcp", ":8084")
//	if err != nil {
//		panic(err)
//	}
//
//	err = rpcServer.Serve(listener)
//	if err != nil {
//		log.Fatal("启动服务出错", err)
//	}
//	fmt.Println("启动grpc服务端成功")
//}
