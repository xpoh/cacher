package cacher

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/xpoh/cacher/pkg/config"
	pb "github.com/xpoh/cacher/pkg/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	cfg   *config.Cfg
	rdb   *redis.Client
	rsync *redsync.Redsync
	pb.UnimplementedCacherServer
}

func (s Server) GetRandomDataStream(in *pb.Request, srv pb.Cacher_GetRandomDataStreamServer) error {
	ctx := context.Background()
	n := len(s.cfg.URLs)
	log.Println("Get request", in.Id)

	var wg sync.WaitGroup
	for i := 0; i < s.cfg.NumberOfRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url := s.cfg.URLs[rand.Intn(n)]
			m := s.rsync.NewMutex(url)
			m.Lock()
			defer m.Unlock()

			var r string
			var err error
			r, err = s.rdb.Get(ctx, url).Result()
			if err == redis.Nil {
				res, err := http.Get(url)
				if err != nil {
					log.Fatal(err)
				}
				b, err := io.ReadAll(res.Body)
				r := string(b)
				res.Body.Close()

				if err != nil {
					log.Fatal(err)
				}

				t := time.Duration(rand.Intn(s.cfg.MaxTimeout-s.cfg.MinTimeout)+s.cfg.MinTimeout) * time.Second
				s.rdb.Set(ctx, url, r, t)
			} else if err != nil {
				panic(err)
			}
			resp := pb.Response{
				Result: r,
			}
			if err := srv.Send(&resp); err != nil {
				log.Printf("send error %v", err)
			}
		}()
	}

	wg.Wait()
	return nil
}

func (s *Server) Run() {
	cfg := config.NewConfig()
	s.cfg = cfg

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
		DB:       0,
	})
	s.rdb = rdb

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Panicln(err)
	}
	pool := goredis.NewPool(rdb)

	rs := redsync.New(pool)
	s.rsync = rs

	// create listiner
	lis, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create grpc Server
	gs := grpc.NewServer()

	pb.RegisterCacherServer(gs, Server{})

	log.Println("start Server")
	// and start...
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
