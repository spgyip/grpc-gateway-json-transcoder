gRPC is broadly used by developing microservices, it's very efficient when both client and server sides are microservices. However, when collaborating with web or native-app, RESTful-JSON API is suitable for end-to-end developing. RESTful-JSON is based on HTTP/JSON which has better ecology for developers debugging/testing.

We can write a RESTful-JSON microservices for each gRPC microservices to translate between the RESTful-JSON/gRPC protocol, but this is awkward. So, in order to communicate to your gRPC service, there are several ways, we are going to demonstrate here.

- [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
- [envoyproxy](https://envoyproxy.io)
- [grpcurl](https://github.com/fullstorydev/grpcurl)

Below features are supported by `Protobuf`/`gRPC`.

- [gRPC transcoding](https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#grpc-transcoding)
- [Protobuf Descriptors](https://buf.build/docs/reference/descriptors)

At the last part, we will introduce another solution [connect-go](https://connectrpc.com/docs/go/getting-started/).

# Topology

This topology shows what we are going to demostrate how to communicate with `greeter_server`

```
         --------------------------------------------------------------------------------------------------------------
         |                                                 greeter_server                                             |
         ---------------------------------------------------(port:50051)-----------------------------------------------
                                                                 /\                                             /reflection
                                                                  |                                                     /\
                                                                [gRPC]                                                   |
                                                                  |                                                    [gRPC]
                    -----------------------------------------------------------------------------------------------      |
                   /\                              /\                               /\                           /\      |
                    |                               |                                |                            |      |
                    |                    -------------------------    ---------------------------                 |      |
                    |                    |      grpc_gateway     |    |         envoy           |                 |      |
                    |                    -------(port:52051)------    --------(port:51051)-------                 |      |
                    |                               /\                            /\       (*helloworld.pb*)      |      |
                    |                                |                             |                              |      |
                    |                                -------------------------------                              |      |
                    |                                          |                                                  |      |
                    |                                      [RESTful]                                              |      |
                    |                                          |                                                  |      |
         -------------------        ---------------------------------------------------------------       ---------------------
         | greeter_client  |        |                     cURL                                    |       |     grpcurl       |
         -------------------        ---------------------------------------------------------------       ---------------------
                                                                                                                         (-protoset=*helloworld.pb*)
```

- Use `greeter_client`, which is a gRPC client, to communicate directly to `greeter_server`.
- `grpc_gateway` as proxy, use `curl` to communicate to proxy with RESTful.
- `envoy` as proxy, use `curl` to communicate to proxy with RESTful. 
    - The `Descriptor` file(*helloworld.pb*) must be provided.
- Use `grpcurl` to communicate directly to `greeter_server`, `grpcurl` supports translation between JSON and protobuf. 
    - The `Descriptor` file(*helloworld.pb*) must be provided.
    - Another option for `Descriptor` is, if the `greeter_server` supports `server reflection`, `grpcurl` will fetch `Descriptor` from server automatically. 


Protocol notation:
- gRPC: http2+proto
- RESTful: http1/http2+JSON

# Step-by-Step

1. Build `greeter` and launch `greeter_server`/`greeter_gateway`

```shell
make bin
```

Open 2 more terminals to run `greeter_server`/`greeter_gateway` seperatly.

```shell
./bin/greeter_server
./bin/greeter_gateway
```

2. Try `greeter_client`

```shell
./bin/greeter_client
```

If all things are right, a response message is responded.

```
2024/02/28 11:24:31 Greeting: Hello world
```

3. Try `curl` with `greeter_gateway`

```shell
sh sbin/restful gateway
```

If all things are right, a response JSON is responded.

```
{"message":"Hello world"}
```

4. Try `grpcurl`

> Notice: helloworld.pb must be provided. 

```shell
sh sbin/grpcurl ./config/envoy/helloworld.pb
```

If all things are right, a response JSON is responded.

```
{
  "message": "Hello hello-from-grpcurl"
}
```

5. Configure and launch envoy

> Because envoy is launched by Docker, you must set your host's localhost IP(it's should be eth1 on most environment) for envoy, allowing it can communicate with `greeter_server` from inside container network.

```shell
sh sbin/configure-envoy 192.168.1.101
```

> `envoy` will be launched by `Docker`, so it's required you are as `root` or having the privilidge for `sudo`.

Open another terminal to run `envoy`.

```shell
sh sbin/runenvoy
```

6. Try RESTful with `envoy`

```shell
sh restful envoy
```

If all things are right, a response JSON is responded.

```
{
 "message": "Hello world"
}
```

# Server reflection

As demonstrations above, there is a bad experience when using `grpcurl` or `envoy` proxy, a `Descriptor`(*helloworld.pb*) file must be provided. We can use [Server reflection](https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md) to make things simpler. The `gRPC` server provides a relection method returning the `Descriptor`, client side can obtain `Descriptor` through this method. Then the *helloworld.pb* file can be omitted at client side.

> It's a pity that server reflection hasn't been supported by envoy yet. 
> 
> [envoy/issue/1182](https://github.com/envoyproxy/envoy/issues/1182).

## Examine the error

You can try with `grpcurl` by removing the *helloworld.pb*, error message will be produced.

```shell
sh sbin/grpcurl
Error invoking method "helloworld.v1.GreeterService/SayHello": failed to query for service descriptor "helloworld.v1.GreeterService": server does not support the reflection API
```

## Use server reflection

Run `greeter_server` with argument `-server_reflection` can turn on server reflection support, can be confirmed by log messsage `Server reflection is on`.

```shell
./bin/greeter_server -server_reflection
2024/02/29 16:18:07 Server reflection is ON.
```

Then try `grpcurl` by removing the *helloworld.pb* again, you can see all things work fine.

```shell
sh sbin/grpcurl
{
  "message": "Hello hello-from-grpcurl"
}
```

# connect-go

Apart from gRPC's official framework, [connect-go](https://connectrpc.com/docs/go/getting-started/) is another framework which is gRPC-compatible, and supports RESTful-JSON as well. 

As in `buf.gen.yaml`, plugin `buf.build/connectrpc/go` is used to generate `connect-go` stub codes, locates at `protogen/helloworld/v1/helloworldv1connect/`. `greeter_server.connect` is implemented at `cmd/connect/greeter_server/`.

Unlike `gRPC` protocol only supports `http2+proto`, `connect-go` protocol supports `http1|http2|http3+proto|JSON` without proxy at all. So, it's compatible with `gRPC` protocol.

## Step-by-step

1. Run `greeter_server.connect`

```shell
./bin/greeter_server.connect
```

2. Try `gRPC client`

`connect-go` is compatible with `gRPC` protocol, so a `gRPC client` can communicate with `connect-go server`.


```shell
./bin/greeter_client -addr "localhost:60051"
```

A successful message should be printed.

```shell
2024/03/06 15:48:33 Greeting: Hello world
```

3. Try `http1+JSON`

```shell
curl -v -X POST --header "Content-Type: application/json" -d '{"name": "helloworld"}' http://127.0.0.1:60051/helloworld.v1.GreeterService/SayHello
```

4. Try `http2+JSON`

> Make sure if your `curl` support `http2`: [curl-http2](https://curl.se/docs/http2.html).
> Specify `--http2-prior-knowledge` to use `http2`.

```shell
curl -v -X POST --http2-prior-knowledge --header "Content-Type: application/json" -d '{"name": "helloworld"}' http://127.0.0.1:60051/helloworld.v1.GreeterService/SayHello
```

5. Try `http1+proto`

> `buf-curl` supports protobuf serialization with `curl`.

```shell
buf curl -v --schema proto --data '{"name": "helloworld"}' http://127.0.0.1:60051/helloworld.v1.GreeterService/SayHello
```

6. Try `http2+proto`

> `buf-curl` supports protobuf serialization with `curl`.
>  Specify `--http2-prior-knowledge` to use `http2`.

```shell
buf curl -v --protocol grpc --http2-prior-knowledge --schema proto --data '{"name": "helloworld"}' http://127.0.0.1:60051/helloworld.v1.GreeterService/SayHello
```
