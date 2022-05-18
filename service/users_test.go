package service

import (
	"context"
	"gRPC_User/helper"
	"gRPC_User/model"
	"gRPC_User/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"os"
	"testing"
	"time"
)

var user = &model.Authentication{
	User:     "admin",
	Password: "admin",
}

func TestProductService_GetUserByName(t *testing.T) {

	testCases := []struct { //定义测试的结构体(测试不同请求)
		TestName   string
		UserDate   string
		StatusCode int32
	}{
		//测试组
		{"Test_Right", "joy", 200},
		{"Test_Unknown", "ttt", 500},
	}

	conn, err := grpc.Dial(":8084", grpc.WithTransportCredentials(helper.GetClientCred()), grpc.WithPerRPCCredentials(user))
	require.NoError(t, err)

	c := proto.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {

			r, err := c.GetUserByName(ctx, &proto.UserGetReq{Name: testCase.UserDate})
			require.NoError(t, err)

			assert.Equal(t, testCase.StatusCode, r.Code, "They should be equal")

		})
	}

}

func TestUserService_GetUsers(t *testing.T) {

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

	conn, err := grpc.Dial(":8084", grpc.WithTransportCredentials(helper.GetClientCred()), grpc.WithPerRPCCredentials(user))
	require.NoError(t, err)

	c := proto.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {

			r, err := c.GetUsers(ctx, &proto.UsersGetReq{
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
		{"Test_Repeat", "joy", 400},
		{"Test_Length", "k", 400},
		{"Test_Right", "jakc", 201},
	}

	conn, err := grpc.Dial(":8084", grpc.WithTransportCredentials(helper.GetClientCred()), grpc.WithPerRPCCredentials(user))

	require.NoError(t, err)
	defer conn.Close()
	c := proto.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	for _, testCase := range testCases { //进行三次测试
		t.Run(testCase.TestName, func(t *testing.T) {
			r, err := c.AddUser(ctx, &proto.UserPostReq{Name: testCase.UserDate})
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
		{"Test_Right", "joy", "joy1", 201},
		{"Test_Length", "1111", "k", 400},
		{"Test_Unknown", "te", "22", 400},
		{"Test_Repeat", "jack", "mark", 400},
		{"Test_Empty", "", "", 400},
	}

	conn, err := grpc.Dial(":8084", grpc.WithTransportCredentials(helper.GetClientCred()), grpc.WithPerRPCCredentials(user))

	require.NoError(t, err)
	defer conn.Close()
	c := proto.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			r, err := c.UpdUserName(ctx, &proto.UserPutReq{
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
		{"Test_Unknown", "want", 400},
		{"Test_Right", "test1111", 201},
	}

	conn, err := grpc.Dial(":8084", grpc.WithTransportCredentials(helper.GetClientCred()), grpc.WithPerRPCCredentials(user))

	require.NoError(t, err)
	defer conn.Close()
	c := proto.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			r, err := c.DelUser(ctx, &proto.UserDelReq{Name: testCase.UserDate})
			require.NoError(t, err)

			assert.Equal(t, testCase.StatusCode, r.Code, "They should be equal")
		})
	}
}

func TestMain(m *testing.M) {

	retCode := m.Run()

	os.Exit(retCode)

}
