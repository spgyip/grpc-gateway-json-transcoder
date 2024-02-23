.PHONY: proto bin clean

bin:
	go build -o bin/ ./...

proto:
	protoc \
		--go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
		--include_imports --include_source_info \
		--descriptor_set_out=helloworld/helloworld.pb \
        helloworld/helloworld.proto

		#--go_opt=Mhelloworld/helloworld.proto=github.com/spgyip/grpc-gateway-json-transcoding/helloworld \
# -o helloworld/helloworld.pb \
#-I ./ \
#-I $(GOPATH)/src/github.com/googleapis/googleapis/  \

clean:
	rm -fv bin/*
