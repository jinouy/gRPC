package main

import (
	"context"
	"gRPC_User/proto/chat"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"testing"
)

var wg sync.WaitGroup

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
	}

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {

			// 创建连接，拨号
			conn, err := grpc.Dial("localhost:9999", grpc.WithTransportCredentials(insecure.NewCredentials()))
			require.NoError(t, err)

			defer conn.Close()

			wg.Add(2)
			go func() {
				client1 := chat.NewOnLineChatClient(conn)
				ctx, cancel := context.WithCancel(context.Background())
				stream1, err := client1.SayHi(ctx)
				require.NoError(t, err)

				defer wg.Done()
				err = stream1.Send(&chat.HiRequest{Name: testCase.UserName1})
				require.NoError(t, err)

				//接收 服务端信息
				reply, err := stream1.Recv()
				require.NoError(t, err)
				if reply.MessageType == 1 {
					return
				}

				var line string
				line = testCase.SayHello1
				if line == "exit" {
					cancel()
					return
				}
				err = stream1.Send(&chat.HiRequest{Message: line})
				require.NoError(t, err)
				if reply.Message != "" {
					return
				} else {
					t.Error("没有传到客户端")
				}
			}()

			go func() {
				// 声明客户端
				client2 := chat.NewOnLineChatClient(conn)
				// 创建双向数据流
				ctx, cancel := context.WithCancel(context.Background())
				stream2, err := client2.SayHi(ctx)
				require.NoError(t, err)

				defer wg.Done()
				err = stream2.Send(&chat.HiRequest{Name: testCase.UserName2})
				require.NoError(t, err)

				//接收 服务端信息
				reply, err := stream2.Recv()
				require.NoError(t, err)
				if reply.MessageType == 1 {
					return
				}

				var line string
				line = testCase.SayHello2
				if line == "exit" {
					cancel()
					return
				}
				err = stream2.Send(&chat.HiRequest{Message: line})
				require.NoError(t, err)
				if reply.Message != "" {
					return
				} else {
					t.Error("没有传到客户端")
				}

			}()
			wg.Wait()

		})
	}
}
