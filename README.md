It's very popular to use RESTful JSON API developing applications, meanwhile we want to write gRPC service only. There are two ways of transcoding from RESTful to gRPC.

- Using [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).
- Using [envoy](https://envoyproxy.io) as proxy.

All these are supported by [gRPC transcoding](https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#grpc-transcoding).

# Notice

This is a demo for how to use gRPC transconding with grpc-gateway and envoy, it's only supported in localhost.

# Topology

```
         ------------------------------------------------------------------------
         |                      greeter_server                                  |
         ----------------(         50051             )---------------------------
                             /\    /\              /\
                              |     |               |
                              |     -----[gRPC]--   --------------[gRPC]---
                              |                 |                         |
                              |       -------------------------    ----------------------
                              |      |   grpc_gateway         |    |     envoy          |
                              |       -------(52051)----------     --------(51051)-------
                              |                 /\                          /\
                    ---[gRPC]--                 |                           |
                    |                       [RESTful]       -----[RESTful] ---
                    |                           |           |
         -------------------        ------------------------------------------
         | greeter_client  |        |                 cURL                   |
         -------------------        ------------------------------------------
```


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
