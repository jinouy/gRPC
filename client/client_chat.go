package main

import (
	"bufio"
	"context"
	"fmt"
	"gRPC_User/client/auth"
	"gRPC_User/comm"
	"gRPC_User/model"
	pb "gRPC_User/proto/chat"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"time"
)

func ConsoleLog(message string) {

	fmt.Printf("\n------ %s -----\n%s\n> ", time.Now().Format("2006-01-02 15:04:05"), message)
}

// 输入
func Input(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	if err != nil {
		if err == io.EOF {
			return "", nil
		} else {
			return "", err
		}
	}
	return string(line), nil
}

func main() {

	token := &model.Auth{
		User: auth.User,
	}
	//  grpc.WithTransportCredentials(comm.GetClientCred())
	// 创建连接，拨号
	conn, err := grpc.Dial(":9998", grpc.WithTransportCredentials(comm.GetClientCred()), grpc.WithPerRPCCredentials(token))
	if err != nil {
		log.Printf("连接失败: [%v] ", err)
		return
	}
	defer conn.Close()

	// 声明客户端
	client := pb.NewOnLineChatClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	//ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("user", token.User))

	// 创建双向数据流
	stream, err := client.SayHi(ctx)
	if err != nil {
		log.Printf("创建数据流失败: [%v] ", err)
		return
	}

	// 接收 服务端信息
	go func() {
		var (
			reply *pb.HiReply
			err   error
		)
		for {
			reply, err = stream.Recv()
			if err != nil {
				log.Printf("数据传输失败: [%v]", err)
				cancel()
				return
			}
			ConsoleLog(reply.Message)

		}
	}()

	go func() {
		var (
			line string
			err  error
		)
		for {
			line, err = Input("")
			if err != nil {
				log.Printf("数据输入失败: [%v]", err)
				return
			}
			if line == "exit" {
				cancel()
				break
			}
			err = stream.Send(&pb.HiRequest{
				Message: line,
			})
			fmt.Print("> ")
			if err != nil {
				log.Printf("数据传输为空: [%v]", err)
				return
			}
		}
	}()

	<-ctx.Done()
	fmt.Println("退出")
}
