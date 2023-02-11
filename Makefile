all: runstack udpechosrv

runstack:
	CGO_ENABLED=0 go build -o bin ./app/runstack

udpechosrv:
	CGO_ENABLED=0 go build -o bin ./app/udpechosrv
