package client

import (
	pb "cacher/pkg/proto"
	"context"
	"io"
	"log"
	"sync"

	"google.golang.org/grpc"
)

type Client struct{}

func (c *Client) Run() {
	// dial cacher
	conn, err := grpc.Dial("server:50005", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with cacher %v", err)
	}
	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := pb.NewCacherClient(conn)
			in := &pb.Request{Id: 1}
			stream, err := client.GetRandomDataStream(context.Background(), in)
			if err != nil {
				log.Println("open stream error %v", err)
				return
			}
			for {
				resp, err := stream.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					log.Fatalf("cannot receive %v", err)
				}
				log.Printf("[client] Resp received: %s", resp.Result)
			}
		}()
	}
	wg.Wait()
	log.Printf("finished")
}
