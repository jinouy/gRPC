package main

import (
	"context"
	"gRPC_User/comm"
	"gRPC_User/model"
	"gRPC_User/proto/chat"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
)

func TestService_SayHi(t *testing.T) {

	testCases := []struct { //定义测试的结构体
		TestName  string
		UserName1 string
		UserName2 string
		SayHello1 string
		SayHello2 string
	}{
		//测试组
		{"Test_Coven", "joy", "jack", "hello", "hi"},
		{"Test_Repeat", "joy", "jack", "exit", "hi"},
		{"Test_Quit", "jack", "jack", "exit", "hi"},
		{"Test_Token", "lily", "jack", "hello", "hi"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {

			token1 := &model.Auth{
				User: testCase.UserName1,
			}

			token2 := &model.Auth{
				User: testCase.UserName2,
			}

			// 创建连接，拨号
			conn1, err := grpc.Dial(":9998", grpc.WithTransportCredentials(comm.GetClientCred()), grpc.WithPerRPCCredentials(token1))
			require.NoError(t, err)

			conn2, err := grpc.Dial(":9998", grpc.WithTransportCredentials(comm.GetClientCred()), grpc.WithPerRPCCredentials(token2))
			require.NoError(t, err)

			defer conn1.Close()
			client1 := chat.NewOnLineChatClient(conn1)
			// 创建双向数据流
			ctx1, cancel1 := context.WithCancel(context.Background())
			//ctx1 = metadata.NewOutgoingContext(ctx1, metadata.Pairs("user", testCase.UserName1))

			stream1, err := client1.SayHi(ctx1)
			require.NoError(t, err)

			client2 := chat.NewOnLineChatClient(conn2)
			// 创建双向数据流
			ctx2, cancel2 := context.WithCancel(context.Background())
			//ctx2 = metadata.NewOutgoingContext(ctx2, metadata.Pairs("user", testCase.UserName2))
			stream2, err := client2.SayHi(ctx2)
			require.NoError(t, err)

			// 接收 服务端信息
			reply, err := stream1.Recv()
			require.NoError(t, err)
			require.Equal(t, "连接成功!", reply.Message)

			reply, err = stream2.Recv()
			require.NoError(t, err)
			require.Equal(t, "连接成功!", reply.Message)

			reply, err = stream1.Recv()
			require.NoError(t, err)
			require.Equal(t, "欢迎 "+testCase.UserName2+"!", reply.Message)

			// 发送信息
			if testCase.SayHello1 != "exit" {

				// 发送 消息
				err = stream1.Send(&chat.HiRequest{Message: testCase.SayHello1})
				require.NoError(t, err)

				reply, err = stream2.Recv()
				require.NoError(t, err)
				require.Equal(t, testCase.UserName1+": "+testCase.SayHello1, reply.Message)

			}

			if testCase.SayHello2 != "exit" {

				//reply, err = stream2.Recv()
				//require.NoError(t, err)
				//require.Equal(t, "欢迎 "+testCase.UserName1+"!", reply.Message)

				err = stream2.Send(&chat.HiRequest{Message: testCase.SayHello2})
				require.NoError(t, err)

				reply, err = stream1.Recv()
				require.NoError(t, err)
				require.Equal(t, testCase.UserName2+": "+testCase.SayHello2, reply.Message)

			}
			if testCase.SayHello1 == "exit" {
				cancel1()
			}
			if testCase.SayHello2 == "exit" {
				cancel2()
			}

		})
	}
}
