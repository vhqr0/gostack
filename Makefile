all: bin/runstack bin/udpechosrv bin/udpechocli

bin/runstack: app/runstack/main.go
	CGO_ENABLED=0 go build -o bin ./app/runstack

bin/udpechosrv: app/udpechosrv/main.go
	CGO_ENABLED=0 go build -o bin ./app/udpechosrv

bin/udpechocli: app/udpechocli/main.go
	CGO_ENABLED=0 go build -o bin ./app/udpechocli
