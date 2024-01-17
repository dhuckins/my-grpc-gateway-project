package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	apiv1 "github.com/dhuckins/my-grpc-gateway-project/gen/proto/go/api/v1"
)

var log = slog.New(slog.NewTextHandler(os.Stderr, nil))

type server struct {
	apiv1.UnimplementedKvServiceServer
	data map[string]string
	lock sync.Locker
}

func (kv *server) Put(_ context.Context, req *apiv1.PutRequest) (*apiv1.PutResponse, error) {
	log.Info(
		"got put request",
		"request.name", req.Name,
	)
	kv.lock.Lock()
	defer kv.lock.Unlock()

	if _, exists := kv.data[req.Name]; exists {
		return nil, status.Errorf(codes.AlreadyExists, "key %q already exists", req.Name)
	}

	kv.data[req.Name] = req.Value

	return &apiv1.PutResponse{}, nil

}

func (kv *server) Get(_ context.Context, req *apiv1.GetRequest) (*apiv1.GetResponse, error) {
	log.Info(
		"got get request",
		"request.name", req.Name,
	)
	kv.lock.Lock()
	defer kv.lock.Unlock()

	value, ok := kv.data[req.Name]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "key %q not found", req.Name)
	}
	return &apiv1.GetResponse{
		Name:  req.Name,
		Value: value,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	apiv1.RegisterKvServiceServer(grpcServer, &server{
		data: map[string]string{},
		lock: &sync.Mutex{},
	})

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}()

	conn, err := grpc.Dial(
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	mux := runtime.NewServeMux()
	if err := apiv1.RegisterKvServiceHandler(context.Background(), mux, conn); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	httpServer := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	log.Info("running gateway at http://localhost:8081")
	if err := httpServer.ListenAndServe(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
