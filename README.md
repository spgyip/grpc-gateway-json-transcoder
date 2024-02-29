gRPC is broadly used by developing microservices, it's very efficient when both client and server side are microservices. However, when you are developing API for web or native-app, RESTful-JSON API is another standard suitable for end-to-end developing. Further, RESTful-JSON which based on HTTP/JSON has better ecology for developers debugging/testing.

We can write a RESTful-JSON microservices for each gRPC microservices to translate between the RESTful-JSON/gRPC protocol, but this is awkward. So, in order to communicate to your gRPC service, there are several ways, we are going to demonstrate here.

- [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).
- [envoyproxy](https://envoyproxy.io).
- [grpcurl](https://github.com/fullstorydev/grpcurl).


Some feature are supported by `Protobuf`/`gRPC`.

- [gRPC transcoding](https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#grpc-transcoding).
- [Protobuf Descriptors](https://buf.build/docs/reference/descriptors)


# Topology

This topology shows that there are several ways to communicate with `greeter_server`

- Use `greeter_client`, which is a gRPC client, to communicate directly to `greeter_server`.
- `grpc_gateway` as proxy, use `curl` to communicate to proxy with RESTful.
- `envoy` as proxy, use `curl` to communicate to proxy with RESTful. 
    - The `Descriptor` file(*helloworld.pb*) must be provided.
- Use `grpcurl` to communicate directly to `greeter_server`, however `grpcurl` supports translate JSON to protobuf. 
    - The `Descriptor` file(*helloworld.pb*) must be provided.
    - Another option for `Descriptor` is, if the `greeter_server` supports `server reflection`, `grpcurl` will fetch `Descriptor` from server automatically. 

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

Protocol notation:
- gRPC: http2+proto
- RESTful: http+JSON

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

As demonstrations above, there is a bad experience when using `grpcurl` or `envoy` proxy, a `Descriptor`(*helloworld.pb*) file must be provided. We can use [Server reflection](https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md) to make things simpler. That's the `gRPC` server provides a method returns the `Descriptor`, the client side can obtain `Descriptor` through the reflection method, which can omit the *helloworld.pb* file at client side.

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
```

