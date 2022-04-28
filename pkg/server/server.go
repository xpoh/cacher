package server

import (
	"context"
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/xpoh/cacher/pkg/config"
	pb "github.com/xpoh/server/pkg/proto"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	cfg   *config.Cfg
	rdb   *goredislib.Client
	rsync *redsync.Redsync
	pb.UnimplementedCacherServer
}

func (s Server) GetRandomDataStream(in *pb.Request, srv pb.Cacher_GetRandomDataStreamServer) error {
	ctx := context.Background()
	n := len(s.cfg.URLs)
	var err error

	var wg sync.WaitGroup
	for i := 0; i < s.cfg.NumberOfRequests; i++ {
		wg.Add(1)
		go func() {
			var r string

			defer wg.Done()
			url := s.cfg.URLs[rand.Intn(n)]

			m := s.rsync.NewMutex("mutex-" + url)
			err = m.Lock()
			if err != nil {
				log.Printf("Error lock mutex %v err=%v", m.Name(), err)
			} else {
				log.Printf("Lock mutex succseful %v", m.Name())
			}
			defer func(m *redsync.Mutex) {
				status, err := m.Unlock()
				log.Printf("Unlock %v status %v", m.Name(), status)
				if err != nil {
					log.Println("Error unlock")
				}
			}(m)

			log.Printf("Find url=%v in redis...", url)

			r, err = s.rdb.Get(ctx, url).Result()
			fmt.Printf("get %v from redis err=%v ", url, err)
			if err == goredislib.Nil {
				log.Printf("NOT FIND")
				res, err := http.Get(url)
				if err != nil {
					log.Printf("Error get request from web %v (%v)", url, err)
					return
				}
				//b, err := io.ReadAll(res.Body)
				//r := string(b)
				r = "from web: " + url

				res.Body.Close()

				if err != nil {
					log.Printf("Error read request body (%v)", err)
				}

				t := time.Duration(rand.Intn(s.cfg.MaxTimeout-s.cfg.MinTimeout)+s.cfg.MinTimeout) * time.Second
				log.Printf("ADD %v to redis with time %v\n", url, t)
				err = s.rdb.Set(ctx, url, url, t).Err()
				if err != nil {
					log.Printf("Error write %v to redis (%v)", url, err)
				}
			} else if err != nil {
				log.Printf("Error Get request from redis %v (%v)\n", url, err)
			} else {
				log.Printf("%v FIND\n", url)
				r = "from cache: " + url
			}

			resp := pb.Response{
				Result: r,
			}
			if err := srv.Send(&resp); err != nil {
				log.Printf("rpc send error (%v)", err)
			}
		}()
	}

	wg.Wait()
	return err
}

func (s *Server) Run() {
	cfg := config.NewConfig()
	s.cfg = cfg

	rdb := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	s.rdb = rdb
	defer func(rdb *goredislib.Client) {
		err := rdb.Close()
		if err != nil {
			log.Println(err)
		}
	}(rdb)

	pool := goredis.NewPool(rdb)
	rs := redsync.New(pool)
	s.rsync = rs

	// create listiner
	lis, err := net.Listen("tcp", "localhost:50005")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create grpc Server
	gs := grpc.NewServer()

	pb.RegisterCacherServer(gs, s)

	log.Println("start Server")
	// and start...
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
