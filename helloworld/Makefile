.PHONY: proto bin clean

bin:
	go build -o bin/ ./...

proto:
	protoc \
		-I ./ \
		-I /usr/local/include/google/googleapis/ \
		--go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
		--include_imports --include_source_info \
		--descriptor_set_out=helloworld/helloworld.pb \
        helloworld/helloworld.proto
	cp -v helloworld/helloworld.pb ../envoy-config/

proto1:
	buf generate -v helloworld
	buf build -o ../envoy-config/helloworld.pb


clean:
	rm -fv bin/*
