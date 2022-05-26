package controller

import (
	"context"
	"fmt"
	"gRPC_User/comm"
	"gRPC_User/model"
	user2 "gRPC_User/proto/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"testing"
	"time"
)

var users = &model.Authentication{
	UserName: "admin",
	Password: "123456",
}

func TestProductService_GetUserByName(t *testing.T) {

	testCases := []struct { //定义测试的结构体(测试不同请求)
		TestName   string
		UserDate   string
		StatusCode int32
	}{
		//测试组
		{"Test_Right", "joy1", 200},
		{"Test_Unknown", "ttt", 500},
	}

	conn, err := grpc.Dial(":8084", grpc.WithTransportCredentials(comm.GetClientCred()), grpc.WithPerRPCCredentials(users))
	require.NoError(t, err)

	c := user2.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {

			r, err := c.GetUserByName(ctx, &user2.UserGetReq{Name: testCase.UserDate})
			require.NoError(t, err)

			assert.Equal(t, testCase.StatusCode, r.Code, "They should be equal")

		})
	}

}

func TestUserService_UserList(t *testing.T) {

	testCases := []struct { //定义测试的结构体
		TestName   string
		NumDate    int32
		SizeDate   int32
		StatusCode int32
	}{
		//测试组
		{"Test_Right", 2, 2, 200},
		{"Test_Num", 100, 1, 400},
		{"Test_Size", 2, 100, 400},
	}

	conn, err := grpc.Dial(":8084", grpc.WithTransportCredentials(comm.GetClientCred()), grpc.WithPerRPCCredentials(users))
	require.NoError(t, err)

	c := user2.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {

			r, err := c.UserList(ctx, &user2.ListGetReq{
				Page:  testCase.NumDate,
				Limit: testCase.SizeDate,
			})
			require.NoError(t, err)

			assert.Equal(t, testCase.StatusCode, r.Code, "They should be equal")

		})
	}

}

func TestUserService_AddUser(t *testing.T) {

	testCases := []struct { //定义测试的结构体
		TestName   string
		UserDate   string
		StatusCode int32
	}{
		//测试组
		{"Test_Repeat", "joy", 500},
		{"Test_Length", "k", 400},
		{"Test_Right", "jakc121", 201},
	}

	conn, err := grpc.Dial(":8084", grpc.WithTransportCredentials(comm.GetClientCred()), grpc.WithPerRPCCredentials(users))

	require.NoError(t, err)
	defer conn.Close()
	c := user2.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	for _, testCase := range testCases { //进行三次测试
		t.Run(testCase.TestName, func(t *testing.T) {
			r, err := c.AddUser(ctx, &user2.UserAddReq{Name: testCase.UserDate})
			require.NoError(t, err)
			assert.Equal(t, testCase.StatusCode, r.Code, "They should be equal")
		})
	}
}

func TestUserService_UpdUserName(t *testing.T) {
	testCases := []struct { //定义测试的结构体
		TestName   string
		OldDate    string
		NewDate    string
		StatusCode int32
	}{
		//测试组
		//{"Test_Right", "joy", "joy11", 201},
		//{"Test_Length", "1111", "k", 400},
		{"Test_Unknown", "te", "22", 500},
		//{"Test_Repeat", "jack", "mark", 500},
	}

	conn, err := grpc.Dial(":8084", grpc.WithTransportCredentials(comm.GetClientCred()), grpc.WithPerRPCCredentials(users))

	require.NoError(t, err)
	defer conn.Close()
	c := user2.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			r, err := c.UpdUserName(ctx, &user2.UserUpdReq{
				OldName: testCase.OldDate,
				NewName: testCase.NewDate,
			})
			require.NoError(t, err)
			assert.Equal(t, testCase.StatusCode, r.Code, "They should be equal")
		})
	}
}

func TestUserService_DelUser(t *testing.T) {
	testCases := []struct { //定义测试的结构体
		TestName   string
		UserDate   string
		StatusCode int32
	}{
		//测试组
		{"Test_Unknown", "11111", 500},
		{"Test_Right", "joy11", 201},
	}

	conn, err := grpc.Dial(":8084", grpc.WithTransportCredentials(comm.GetClientCred()), grpc.WithPerRPCCredentials(users))

	require.NoError(t, err)
	defer conn.Close()
	c := user2.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			r, err := c.DelUser(ctx, &user2.UserDelReq{Name: testCase.UserDate})
			require.NoError(t, err)

			assert.Equal(t, testCase.StatusCode, r.Code, "They should be equal")
		})
	}
}

func TestMain(m *testing.M) {
	go InitServer()
	time.Sleep(2 * time.Second)
	res := m.Run()
	os.Exit(res)
}

func InitServer() {

	rpcServer := grpc.NewServer(grpc.Creds(comm.GetCertService()), grpc.UnaryInterceptor(comm.GetToke()))

	src := &UserService{}
	user2.RegisterUserServiceServer(rpcServer, src)

	listener, err := net.Listen("tcp", ":8084")
	if err != nil {
		panic(err)
	}

	err = rpcServer.Serve(listener)
	if err != nil {
		log.Fatal("启动服务出错", err)
	}
	fmt.Println("启动grpc服务端成功")

}
