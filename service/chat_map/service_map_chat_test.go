package main

import (
	"context"
	"gRPC_User/client/auth"
	"gRPC_User/model"
	"gRPC_User/proto/chat"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		//{"Test_Coven", "joy", "jack", "hello", "hi"},
		{"Test_Repeat", "joy", "jack", "exit", "hi"},
		//{"Test_Quit", "jack", "jack", "exit", "hi"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {

			var err error
			var opts []grpc.DialOption

			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
			opts = append(opts, grpc.WithPerRPCCredentials(new(model.Auth)))
			opts = append(opts, grpc.WithUnaryInterceptor(auth.Clientinerceptor))

			// 创建连接，拨号
			conn, err := grpc.Dial("localhost:9999", opts...)
			require.NoError(t, err)

			defer conn.Close()
			client1 := chat.NewOnLineChatClient(conn)
			ctx1, cancel1 := context.WithCancel(context.Background())

			stream1, err := client1.SayHi(ctx1)
			require.NoError(t, err)

			client2 := chat.NewOnLineChatClient(conn)
			// 创建双向数据流
			ctx2, cancel2 := context.WithCancel(context.Background())
			stream2, err := client2.SayHi(ctx2)
			require.NoError(t, err)

			// 发送信息
			//err = stream1.Send(&chat.HiRequest{Name: testCase.UserName1})
			//require.NoError(t, err)
			if testCase.SayHello1 != "exit" {
				err = stream1.Send(&chat.HiRequest{Message: testCase.SayHello1})
				require.NoError(t, err)

				//接收 服务端信息
				reply, err := stream1.Recv()
				require.NoError(t, err)
				//if reply.MessageType == 1 {
				//	return
				//}
				require.Equal(t, reply.Message, testCase.UserName2+": "+testCase.SayHello2)
			}

			//err = stream2.Send(&chat.HiRequest{Name: testCase.UserName2})
			//require.NoError(t, err)
			if testCase.SayHello2 != "exit" {
				err = stream2.Send(&chat.HiRequest{Message: testCase.SayHello2})
				require.NoError(t, err)

				//接收 服务端信息
				reply, err := stream2.Recv()
				require.NoError(t, err)
				//if reply.MessageType == 1 {
				//	return
				//}
				require.Equal(t, reply.Message, testCase.UserName1+": "+testCase.SayHello1)
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
