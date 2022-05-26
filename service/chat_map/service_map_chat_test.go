package main

import (
	"context"
	"gRPC_User/comm"
	"gRPC_User/model"
	"gRPC_User/proto/chat"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"testing"
)

func TestService_SayHi(t *testing.T) {

	var mutex sync.Mutex

	testCaseByOne := []struct { //定义测试的结构体：结构体user为单个
		TestName string
		UserName string
		SayHello string
		wantErr  error
	}{
		{"TestUnAuthUser", "", "", status.Errorf(codes.Unauthenticated, "token 不合法")},
		{"TestAuthUser", "joy", "hi", nil},
	}

	for _, testcaseByOne := range testCaseByOne {
		t.Run(testcaseByOne.TestName, func(t *testing.T) {

			token := &model.Auth{
				User: testcaseByOne.UserName,
			}

			conn, err := grpc.Dial(":9998", grpc.WithTransportCredentials(comm.GetClientCred()), grpc.WithPerRPCCredentials(token))
			require.NoError(t, err)

			defer conn.Close()

			client := chat.NewOnLineChatClient(conn)
			ctx, _ := context.WithCancel(context.Background())

			stream, err := client.SayHi(ctx)
			require.NoError(t, err)

			// 接收连接成功的消息
			reply, err := stream.Recv()
			if testcaseByOne.wantErr != nil {
				require.Equal(t, testcaseByOne.wantErr, err)
				return
			}

			require.Equal(t, "连接成功!", reply.Message)

			// 发送消息
			err = stream.Send(&chat.HiRequest{Message: testcaseByOne.SayHello})
			require.NoError(t, err)

		})
	}

	testCases := []struct { //定义测试的结构体
		TestName  string
		UserName1 string
		UserName2 string
		SayHello1 string
		SayHello2 string
		wantErr   error
	}{
		//测试组
		{"TestMultipleAuthUserChat", "joy", "jack", "hello", "hi", nil},
		{"TestMultipleAuthUserLeave", "joy", "jack", "exit", "hi", nil},
		{"TestMultipleAuthSame", "jack", "jack", "exit", "hi", status.Errorf(codes.Unimplemented, "名字已经存在")},
	}

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {

			mutex.Lock()

			defer mutex.Unlock()

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
			if testCase.wantErr != nil {
				require.Equal(t, testCase.wantErr, err)
				return
			}
			require.Equal(t, "连接成功!", reply.Message)

			reply, err = stream1.Recv()
			require.NoError(t, err)
			require.Equal(t, "欢迎 "+testCase.UserName2+"!", reply.Message)

			//reply, err = stream2.Recv()
			//require.NoError(t, err)
			//require.Equal(t, "欢迎 "+testCase.UserName1+"!", reply.Message)

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
				reply, err = stream2.Recv()
				require.NoError(t, err)
				require.Equal(t, testCase.UserName1+" 离开房间", reply.Message)
			}
			if testCase.SayHello2 == "exit" {
				cancel2()
				reply, err = stream1.Recv()
				require.NoError(t, err)
				require.Equal(t, testCase.UserName2+" 离开房间", reply.Message)
			}

		})
	}

	testCaseByMore := []struct { //定义测试的结构体：结构体user为多个
		TestName  string
		UserName1 string
		UserName2 string
		UserName3 string
		SayHello1 string
		SayHello2 string
		SayHello3 string
	}{
		{"TestOrderedMessage", "joy", "jack", "tom", "hello", "hi", "okok"},
	}

	for _, testcaseByMore := range testCaseByMore {
		t.Run(testcaseByMore.TestName, func(t *testing.T) {

			mutex.Lock()

			defer mutex.Unlock()

			token1 := &model.Auth{
				User: testcaseByMore.UserName1,
			}

			token2 := &model.Auth{
				User: testcaseByMore.UserName2,
			}

			token3 := &model.Auth{
				User: testcaseByMore.UserName3,
			}

			conn1, err := grpc.Dial(":9998", grpc.WithTransportCredentials(comm.GetClientCred()), grpc.WithPerRPCCredentials(token1))
			require.NoError(t, err)

			conn2, err := grpc.Dial(":9998", grpc.WithTransportCredentials(comm.GetClientCred()), grpc.WithPerRPCCredentials(token2))
			require.NoError(t, err)

			conn3, err := grpc.Dial(":9998", grpc.WithTransportCredentials(comm.GetClientCred()), grpc.WithPerRPCCredentials(token3))
			require.NoError(t, err)

			defer conn1.Close()
			defer conn2.Close()
			defer conn3.Close()

			client1 := chat.NewOnLineChatClient(conn1)
			ctx1, _ := context.WithCancel(context.Background())

			client2 := chat.NewOnLineChatClient(conn2)
			ctx2, _ := context.WithCancel(context.Background())

			client3 := chat.NewOnLineChatClient(conn3)
			ctx3, _ := context.WithCancel(context.Background())

			stream1, err := client1.SayHi(ctx1)
			require.NoError(t, err)

			stream2, err := client2.SayHi(ctx2)
			require.NoError(t, err)

			stream3, err := client3.SayHi(ctx3)
			require.NoError(t, err)

			// 接收连接成功的消息
			reply, err := stream1.Recv()
			require.NoError(t, err)
			require.Equal(t, "连接成功!", reply.Message)

			reply, err = stream2.Recv()
			require.NoError(t, err)
			require.Equal(t, "连接成功!", reply.Message)

			reply, err = stream3.Recv()
			require.NoError(t, err)
			require.Equal(t, "连接成功!", reply.Message)

			reply, err = stream1.Recv()
			require.NoError(t, err)
			require.Equal(t, "欢迎 "+testcaseByMore.UserName2+"!", reply.Message)

			reply, err = stream1.Recv()
			require.NoError(t, err)
			require.Equal(t, "欢迎 "+testcaseByMore.UserName3+"!", reply.Message)

			reply, err = stream2.Recv()
			require.NoError(t, err)
			require.Equal(t, "欢迎 "+testcaseByMore.UserName3+"!", reply.Message)

			// 发送消息
			err = stream1.Send(&chat.HiRequest{Message: testcaseByMore.SayHello1})
			require.NoError(t, err)

			reply, err = stream2.Recv()
			require.NoError(t, err)
			require.Equal(t, testcaseByMore.UserName1+": "+testcaseByMore.SayHello1, reply.Message)

			reply, err = stream3.Recv()
			require.NoError(t, err)
			require.Equal(t, testcaseByMore.UserName1+": "+testcaseByMore.SayHello1, reply.Message)

			err = stream2.Send(&chat.HiRequest{Message: testcaseByMore.SayHello2})
			require.NoError(t, err)

			reply, err = stream1.Recv()
			require.NoError(t, err)
			require.Equal(t, testcaseByMore.UserName2+": "+testcaseByMore.SayHello2, reply.Message)

			reply, err = stream3.Recv()
			require.NoError(t, err)
			require.Equal(t, testcaseByMore.UserName2+": "+testcaseByMore.SayHello2, reply.Message)

			err = stream3.Send(&chat.HiRequest{Message: testcaseByMore.SayHello3})
			require.NoError(t, err)

			reply, err = stream1.Recv()
			require.NoError(t, err)
			require.Equal(t, testcaseByMore.UserName3+": "+testcaseByMore.SayHello3, reply.Message)

			reply, err = stream2.Recv()
			require.NoError(t, err)
			require.Equal(t, testcaseByMore.UserName3+": "+testcaseByMore.SayHello3, reply.Message)
		})
	}
}
