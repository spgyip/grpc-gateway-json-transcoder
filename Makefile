.PHONY: bin proto clean check-buf-install

bin:
	go build -o bin/ ./...

proto: check-buf-install
	buf lint proto
	buf generate proto
	buf build proto -o config/envoy/helloworld.pb

clean:
	rm -fv bin/*
