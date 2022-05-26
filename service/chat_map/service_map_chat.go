package main

import (
	"fmt"
	"gRPC_User/comm"
	pb "gRPC_User/proto/chat"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"sync"
	"time"
)

type ConnectPool struct {
	Map   map[string]pb.OnLineChat_SayHiServer
	Lock  *sync.RWMutex
	mutex sync.Mutex
}

func NewConcurMap() *ConnectPool {
	return &ConnectPool{
		Map:  make(map[string]pb.OnLineChat_SayHiServer, 10),
		Lock: &sync.RWMutex{},
	}
}

func (p *ConnectPool) Get(name string) pb.OnLineChat_SayHiServer {
	p.Lock.RLock()

	defer p.Lock.RUnlock()

	if stream, ok := p.Map[name]; ok {
		return stream
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

	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Lock.RLock()
	defer p.Lock.RUnlock()

	log.Printf("message: %s\n", message)
	for username, stream_i := range p.Map {
		stream := stream_i
		if username != from {
			err := stream.Send(&pb.HiReply{
				Message: message,
				TS:      &timestamp.Timestamp{Seconds: time.Now().Unix()},
			})
			if err != nil {
				return false
			}
		}
	}
	return true
}

type Service struct {
	pb.UnimplementedOnLineChatServer
}

var connect_pool *ConnectPool

func (s *Service) SayHi(stream pb.OnLineChat_SayHiServer) error {
	//return errors.New("test err")
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.Unimplemented, "no token")
	}
	//for k, v := range md {
	//	fmt.Printf("key:%s value:%s\n", k, v)
	//}
	username := md["user"][0]

	if connect_pool.Get(username) != nil {
		return status.Errorf(codes.Unimplemented, "名字已经存在")
	} else { // 连接成功
		connect_pool.Add(username, stream)
		err := stream.Send(&pb.HiReply{
			Message: fmt.Sprintf("连接成功!"),
		})
		if err != nil {
			return err
		}
	}

	go func() {
		// 阻塞住，等待断开连接的时候触发
		<-stream.Context().Done()
		connect_pool.Del(username)
		connect_pool.BroadCast(username, fmt.Sprintf("%s 离开房间", username))
	}()

	//// 广播，xxxx进入了聊天室直播间
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
	// grpc.Creds(comm.GetCertService())
	// 实例化grpc Server，并开启拦截器
	ser := grpc.NewServer(grpc.Creds(comm.GetCertService()), grpc.StreamInterceptor(comm.GetServerInterceptor()))
	//ser := grpc.NewServer()
	pb.RegisterOnLineChatServer(ser, &Service{}) //必须实现protoes中定义的方法，不然这里无法通过检测

	connect_pool = NewConcurMap()
	// 监听一个 地址:端口
	address, err := net.Listen("tcp", ":9998")
	if err != nil {
		log.Printf("Failed to listen: [%v]", err)
		return
	}
	// 启动服务
	if err := ser.Serve(address); err != nil {
		log.Printf("Failed to start: [%v]", err)
		return
	}
}
