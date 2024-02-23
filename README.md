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

`sh configure-localhost 1.1.1.1`.

2. Build helloworld

```shell
cd helloworld
make bin
```

3. Launch `greeter_server`/`greeter_gateway`

Open 2 different terminals to run these 2 servers.

```shell
cd helloworld
./bin/greeter_server
./bin/greeter_gateway
```

Open another terminal, to check with `greeter_client` that `greeter_server` is ok.

```shell
./bin/greeter_client
```

4. Try RESTful with `greeter_gateway`

```shell
sh restful gateway
```

5. Launch envoy

```shell
sh runenvoy
```

6. Try RESTful with `envoy`

```shell
sh restful envoy
```
