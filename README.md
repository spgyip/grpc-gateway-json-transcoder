RESTful JSON API is very broadly used to develope web applications, also it's really efficient to test with toolkits like curl. Meanwhile we want to write gRPC service only. So, it's realistic to support RESTful communication to gRPC server. Let's demonstrate several ways.

- [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).
- [envoyproxy](https://envoyproxy.io).
- [grpcurl](https://github.com/fullstorydev/grpcurl).

Few features are supported by `gRPC`/`protobuf`.

- [gRPC transcoding](https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#grpc-transcoding).
- [Descriptors](https://buf.build/docs/reference/descriptors)

> This demo is simple, and should be runned as local.

# Topology

```
         --------------------------------------------------------------------------------------------------------------
         |                                                 greeter_server                                             |
         ---------------------------------------------------(port:50051)-----------------------------------------------
                                                                 /\              
                                                                  |              
                                                                [gRPC]
                                                                  |
                    ------------------------------------------------------------------------------------------------------
                   /\                              /\                               /\                            |      |
                    |                               |                                |                            |      |
                    |                    -------------------------    ---------------------------                 |      |
                    |                    |      grpc_gateway     |    |         envoy           |                 |      |
                    |                    -------(port:52051)------    --------(port:51051)-------                 |   /reflection 
                    |                               /\                            /\                              |      |
                    |                                |                             |                              |      |  |*helloworld.pb*|
                    |                                -------------------------------                              |      |        |
                    |                                          |                                                  |      |        |
                    |                                      [RESTful]                                              |      |     (-protoset)
                    |                                          |                                                  |      |        |
         -------------------        ---------------------------------------------------------------       ---------------------   |
         | greeter_client  |        |                     cURL                                    |       |     grpcurl       |<--|
         -------------------        ---------------------------------------------------------------       ---------------------
```

Protocol explanation:
- gRPC: http2+proto, no TLS default.
- RESTful: http+JSON


# How-to

1. Configure localhost ip

> Because envoy runs with docker, we must set localhost IP for envoy to allow communicating with `greeter_server` from inside container.

`sh sbin/configure-envoy 1.1.1.1`.

2. Build `greeter` and launch `greeter_server`/`greeter_gateway`

```shell
make bin
```

Open 2 more terminals to run `greeter_server`/`greeter_gateway`

```shell
./bin/greeter_server
./bin/greeter_gateway
```

Open 1 more terminal, run `greeter_client` to check `greeter_server` is ok.

```shell
./bin/greeter_client
```

If all things are right, a response message is responded.

```
2024/02/28 11:24:31 Greeting: Hello world
```

3. Try RESTful with `greeter_gateway`

```shell
sh sbin/restful gateway
```

If all things are right, a response JSON is responded.

```
{"message":"Hello world"}
```

4. Launch envoy

> `envoy` will be launched with `Docker`, thus it's required you are as `root` or having the privilidge for `sudo`.

Open another terminal to run `envoy`.

```shell
sh sbin/runenvoy
```

5. Try RESTful with `envoy`

```shell
sh restful envoy
```

If all things are right, a response JSON is responded.

```
{
 "message": "Hello world"
}
```
