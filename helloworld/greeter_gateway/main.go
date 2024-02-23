package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gw "github.com/spgyip/grpc-gateway-json-transconding/helloworld/helloworld"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:50051", "Upstream gRPC server endpoint")
	listenAddr         = flag.String("listen", ":52051", "Listen address")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterGreeterHandlerFromEndpoint(
		ctx,
		mux,
		*grpcServerEndpoint,
		opts,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.Printf("Listening gateway %v\n", *listenAddr)
	log.Printf("Upstream gRPC endpoint %v\n", *grpcServerEndpoint)
	err = http.ListenAndServe(*listenAddr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
