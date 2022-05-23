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

const (
	HiReply_CONNECT_SUCCESS chat.HiReply_MessageType = 0
	HiReply_CONNECT_FAILED  chat.HiReply_MessageType = 1
	HiReply_NORMAL_MESSAGE  chat.HiReply_MessageType = 2
)

func TestService_SayHi(t *testing.T) {

	testCases := []struct { //定义测试的结构体
		TestName    string
		UserName    string
		SayHi       string
		MessageType chat.HiReply_MessageType
	}{
		//测试组
		{"Test1", "joy", "1111", HiReply_CONNECT_SUCCESS},
		{"Test2", "k", "1111", HiReply_CONNECT_SUCCESS},
		{"Test3", "jakc121", "1111", HiReply_CONNECT_FAILED},
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

			user := &chat.HiRequest{Name: testCase.UserName}

			err = stream.Send(user)

			//connected := make(chan bool)

			go func() {
				var (
					reply *chat.HiReply
					err   error
				)
				for {
					reply, err = stream.Recv()
					require.NoError(t, err)
				}
				assert.Equal(t, testCase.MessageType, reply.MessageType, "They should be equal")

			}()

		})
	}

}
