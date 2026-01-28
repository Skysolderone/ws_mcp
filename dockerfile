from golang:1.25.6-alpine

workdir /app

copy go.mod go.sum ./

run go mod download

copy . .

run CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o mcp-server cmd/trade/trade.go


cmd ["/app/mcp-server"]