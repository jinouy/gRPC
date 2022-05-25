package main

import (
	"bufio"
	"context"
	"fmt"
	"gRPC_User/client/auth"
	"gRPC_User/model"
	pb "gRPC_User/proto/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var mutex sync.Mutex

// 这是一个加锁的输出，防止乱序或中间插入print数据
func ConsoleLog(message string) {
	mutex.Lock()
	defer mutex.Unlock()
	t := time.Now()
	fmt.Printf("\n------ %s -----\n%s\n> ", t.UTC().Format("2006-01-02 15:04:05"), message)
}

// 输入
func Input(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	if err != nil {
		if err == io.EOF {
			return ""
		} else {
			panic(err)
		}
	}
	return string(line)
}

func main() {

	var err error
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	opts = append(opts, grpc.WithPerRPCCredentials(new(model.Auth)))

	opts = append(opts, grpc.WithUnaryInterceptor(auth.Clientinerceptor))

	//var Users = &model.Auth{
	//	User: auth.InputName(),
	//}
	// 创建连接，拨号
	conn, err := grpc.Dial("localhost:9999", opts...)
	if err != nil {
		log.Printf("连接失败: [%v] ", err)
		return
	}
	defer conn.Close()

	// 声明客户端
	client := pb.NewOnLineChatClient(conn)

	ctx, cancel := context.WithCancel(context.Background())

	// 创建双向数据流
	stream, err := client.SayHi(ctx)
	if err != nil {
		log.Printf("创建数据流失败: [%v] ", err)
	}
	//user := &pb.HiRequest{Name: Users.User}
	//err = stream.Send(user)

	// 创建了一个连接管道
	//connected := make(chan bool)

	// 接收 服务端信息
	go func() {
		var (
			reply *pb.HiReply
			err   error
		)
		for {
			reply, err = stream.Recv()
			if err != nil {
				panic(err)
			}
			ConsoleLog(reply.Message)

			//if reply.MessageType == pb.HiReply_CONNECT_FAILED { // code=1 连接失败
			//	cancel()
			//	break
			//}
			//if reply.MessageType == pb.HiReply_CONNECT_SUCCESS { // code=0 连接成功
			//	connected <- true
			//}
			// 基本都是两个if都不执行，去下一次循环,返回的是 code=2 正常消息
		}
	}()

	go func() {
		//<-connected
		var (
			line string
			err  error
		)
		for {
			line = Input("")
			if line == "exit" {
				cancel()
				break
			}
			err = stream.Send(&pb.HiRequest{
				Message: line,
			})
			fmt.Print("> ")
			if err != nil {
				panic(err)
			}
		}
	}()

	<-ctx.Done()
	fmt.Println("退出")
}
