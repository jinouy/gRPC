package main

import (
	"fmt"
	pb "gRPC_User/proto/chat"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

type ConnectPool struct {
	Map  map[string]interface{}
	Lock *sync.RWMutex
}

func NewConcurMap() *ConnectPool {
	return &ConnectPool{
		Map:  make(map[string]interface{}, 10),
		Lock: &sync.RWMutex{},
	}
}

func (p *ConnectPool) Get(name string) pb.OnLineChat_SayHiServer {
	p.Lock.RLock()

	defer p.Lock.RUnlock()

	if stream, ok := p.Map[name]; ok {
		return stream.(pb.OnLineChat_SayHiServer)
	} else {
		return nil
	}
}

func (p *ConnectPool) Add(name string, stream pb.OnLineChat_SayHiServer) {
	p.Lock.Lock()

	defer p.Lock.Unlock()

	p.Map[name] = stream.(pb.OnLineChat_SayHiServer)
}

func (p *ConnectPool) Del(name string) {
	p.Lock.Lock()

	defer p.Lock.Unlock()

	delete(p.Map, name)
}

func (p *ConnectPool) BroadCast(from, message string) bool {
	log.Printf("message: %s\n", message)
	for username, stream_i := range p.Map {

		stream := stream_i.(pb.OnLineChat_SayHiServer)

		if username != from {
			stream.Send(&pb.HiReply{
				Message:     message,
				MessageType: pb.HiReply_NORMAL_MESSAGE, // 2.正常数据
				TS:          &timestamp.Timestamp{Seconds: time.Now().Unix()},
			})
		}
	}
	return true
}

type Service struct {
	pb.UnimplementedOnLineChatServer
}

var connect_pool *ConnectPool

func (s *Service) SayHi(stream pb.OnLineChat_SayHiServer) error {

	//connect_pool = NewConcurMap()
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

}

func main() {

	connect_pool = NewConcurMap()
	// 监听一个 地址:端口
	address, err := net.Listen("tcp", ":9999")
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
