package main

import (
	"context"
	"gRPC_User/proto/chat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func TestService_SayHi(t *testing.T) {

	// 创建连接，拨号
	conn, err := grpc.Dial("localhost:9999", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	defer conn.Close()

	testCases := []struct { //定义测试的结构体
		TestName    string
		UserName    string
		SayHi       string
		MessageType chat.HiReply_MessageType
	}{
		//测试组
		{"Test_Right", "joy", "hello", 0},
		{"Test_Repeat", "joy", "hi", 1},
		{"Test_Exit", "jack", "exit", 0},
		{"Test_Inner", "jack", "你好", 0},
	}

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {

			// 声明客户端
			client := chat.NewOnLineChatClient(conn)

			ctx, cancel := context.WithCancel(context.Background())

			// 创建双向数据流
			stream, err := client.SayHi(ctx)
			require.NoError(t, err)

			err = stream.Send(&chat.HiRequest{Name: testCase.UserName})
			require.NoError(t, err)

			//接收 服务端信息
			reply, err := stream.Recv()
			require.NoError(t, err)

			assert.Equal(t, testCase.MessageType, reply.MessageType, "The conn should be equal")
			if reply.MessageType == 1 {
				return
			}

			var line string
			line = testCase.SayHi
			if line == "exit" {
				cancel()
				return
			}
			err = stream.Send(&chat.HiRequest{Message: line})
			require.NoError(t, err)

		})
	}
}
