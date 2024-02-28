.PHONY: bin proto clean

bin:
	go build -o bin/ ./...

proto:
	buf lint proto
	buf generate proto
	buf build proto -o config/envoy/helloworld.pb

clean:
	rm -fv bin/*
