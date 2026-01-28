run:
	go run cmd/rsi.go
build:
	GOOS=linux GOARCH=amd64 go build -o rsi cmd/rsi.go
scp-wws:build
	scp rsi wws:~/rsi
	rm rsi