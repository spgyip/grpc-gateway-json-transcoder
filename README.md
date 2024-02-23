It's very popular to use RESTful JSON API developing applications, meanwhile we want to write gRPC service only. There are two ways of transcoding from RESTful to gRPC.

- Using [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).
- Using [envoy](https://envoyproxy.io) as proxy.

All these are supported by [gRPC transcoding](https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#grpc-transcoding).
