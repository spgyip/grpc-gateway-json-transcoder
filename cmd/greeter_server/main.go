/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	helloworldv1 "github.com/spgyip/grpc-gateway-json-transconding/protogen/helloworld/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port              = flag.Int("port", 50051, "The server port")
	server_reflection = flag.Bool("server_reflection", false, "Support server reflection")
)

// server is used to implement helloworld.GreeterServiceServer.
type server struct {
	helloworldv1.UnimplementedGreeterServiceServer
}

// SayHello implements helloworld.GreeterServiceServer
func (s *server) SayHello(ctx context.Context, in *helloworldv1.SayHelloRequest) (*helloworldv1.SayHelloResponse, error) {
	log.Println("Received: ", in)
	return &helloworldv1.SayHelloResponse{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	helloworldv1.RegisterGreeterServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	if *server_reflection {
		log.Println("Server reflection is ON.")
		reflection.Register(s)
	} else {
		log.Println("Server reflection is OFF.")
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
