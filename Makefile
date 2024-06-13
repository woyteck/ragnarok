gateway:
	@go build -C gateway -o ../bin/gateway .
	@./bin/gateway

.PHONY: gateway
