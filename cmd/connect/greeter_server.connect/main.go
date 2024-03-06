package main

import (
	"fmt"
	"net/http"
	"context"
	"log"

	helloworldv1connect "github.com/spgyip/grpc-gateway-json-transcoding/protogen/helloworld/v1/helloworldv1connect"
	helloworldv1 "github.com/spgyip/grpc-gateway-json-transcoding/protogen/helloworld/v1"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"connectrpc.com/connect"
)

const address = ":60051"

func main() {
	mux := http.NewServeMux()
	path, handler := helloworldv1connect.NewGreeterServiceHandler(&serverImpl{})
	mux.Handle(path, handler)
	fmt.Println("Listening on", address)
	http.ListenAndServe(
		address,
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

type serverImpl struct {
	helloworldv1connect.UnimplementedGreeterServiceHandler
}


func (s *serverImpl) SayHello(ctx context.Context, req *connect.Request[helloworldv1.SayHelloRequest]) (*connect.Response[helloworldv1.SayHelloResponse], error) {
  log.Println("Received: ", req.Msg)
  return connect.NewResponse(&helloworldv1.SayHelloResponse{ Message: "Hello " + req.Msg.GetName()}), nil
}
