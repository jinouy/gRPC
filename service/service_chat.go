package main

import (
	"fmt"
	pb "gRPC_User/proto/chat"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// 定义一个类，继承字典(异步,带锁的),一会存入grpc stream对象 { name : stream<obj> }
type ConnectPool struct {
	sync.Map
}

//为这个 类<对象池> 添加方法，分别为 Get,Add,Del和BroadCast(广播信息,群发)
func (p *ConnectPool) Get(name string) pb.OnLineChat_SayHiServer {
	if stream, ok := p.Load(name); ok {
		return stream.(pb.OnLineChat_SayHiServer)
	} else {
		return nil
	}
}

func (p *ConnectPool) Add(name string, stream pb.OnLineChat_SayHiServer) {
	p.Store(name, stream)
}

func (p *ConnectPool) Del(name string) {
	p.Delete(name)
}

// 聊天室内 广播
func (p *ConnectPool) BroadCast(from, message string) {
	log.Printf("message: %s\n", message)
	p.Range(func(username_i, stream_i interface{}) bool {
		username := username_i.(string)
		stream := stream_i.(pb.OnLineChat_SayHiServer)
		if username == from {
			return true
		} else {
			stream.Send(&pb.HiReply{
				Message:     message,
				MessageType: pb.HiReply_NORMAL_MESSAGE, // 2.正常数据
				TS:          &timestamp.Timestamp{Seconds: time.Now().Unix()},
			})
		}
		return true
	})
}

var connect_pool *ConnectPool

// 定义服务器类
type Service struct {
	pb.UnimplementedOnLineChatServer
}

func (s *Service) SayHi(stream pb.OnLineChat_SayHiServer) error {

	//md, _ := metadata.FromIncomingContext(stream.Context())
	//username := md["name"][0] // 从metadata中获取用户名信息，可以理解为请求头里的数据
	recv, err := stream.Recv()
	if err != nil {
		return nil
	}
	username := recv.Name
	if connect_pool.Get(username) != nil {
		stream.Send(&pb.HiReply{
			Message:     fmt.Sprintf("名字: %s 已经存在", username),
			MessageType: pb.HiReply_CONNECT_FAILED, // 1. 连接失败 ， 重名了 用户已经存在
		})
		return nil

	} else { // 连接成功
		connect_pool.Add(username, stream)
		stream.Send(&pb.HiReply{
			Message:     fmt.Sprintf("连接成功!"),
			MessageType: pb.HiReply_CONNECT_SUCCESS, // 0 连接成功
		})
	}

	go func() {
		// 阻塞住，等待断开连接的时候触发
		<-stream.Context().Done()
		connect_pool.Del(username)
		connect_pool.BroadCast(username, fmt.Sprintf("%s 离开房间", username))
	}()

	// 广播，xxxx进入了聊天室直播间
	connect_pool.BroadCast(username, fmt.Sprintf("欢迎 %s!", username))

	//  阻塞接收 该用户后续传来的消息
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		connect_pool.BroadCast(username, fmt.Sprintf("%s: %s", username, req.Message))
	}
	return nil
}

func GetListen() string {
	if len(os.Args) < 2 {
		return ":9999"
	}
	return os.Args[1]
}

func main() {
	connect_pool = &ConnectPool{}

	// 监听一个 地址:端口
	address, err := net.Listen("tcp", GetListen())
	if err != nil {
		panic(err)
	}

	// 注册 grpc
	ser := grpc.NewServer()
	pb.RegisterOnLineChatServer(ser, &Service{}) //必须实现protoes中定义的方法，不然这里无法通过检测

	// 启动服务
	if err := ser.Serve(address); err != nil {
		panic(err)
	}

}
