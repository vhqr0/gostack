all: runstack udpechosrv

runstack:
	CGO_ENABLED=0 go build -o bin ./cmd/runstack

udpechosrv:
	CGO_ENABLED=0 go build -o bin/app ./app/udpechosrv
