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

	testCases := []struct { //定义测试的结构体
		TestName    string
		UserName    string
		SayHi       string
		MessageType chat.HiReply_MessageType
	}{
		//测试组
		{"Test1", "joy", "hello", 0},
		{"Test2", "joy", "hi", 1},
		{"Test3", "jack", "bye", 0},
	}

	// 创建连接
	grpcConn, err := grpc.Dial(":9999", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	// 关闭连接
	defer grpcConn.Close()

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			// 声明客户端
			chatTest := chat.NewOnLineChatClient(grpcConn)

			// 声明context
			ctx, _ := context.WithCancel(context.Background())

			// 创建双向数据流
			stream, err := chatTest.SayHi(ctx)
			require.NoError(t, err)

			go func() {
				err = stream.Send(&chat.HiRequest{
					Name:    testCase.UserName,
					Message: testCase.SayHi,
				})
				require.NoError(t, err)
			}()

			//connected := make(chan bool)

			go func() {
				var (
					reply *chat.HiReply
					err   error
				)
				for {
					reply, err = stream.Recv()
					require.NoError(t, err)
					assert.Equal(t, testCase.MessageType, reply.MessageType, "They should be equal")
				}
			}()

		})
	}

}
