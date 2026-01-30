#!/bin/bash
# 生成proto Go代码

set -e

PROJECT_ROOT=$(cd "$(dirname "$0")/.." && pwd)
PROTO_DIR="$PROJECT_ROOT/proto"
OUT_DIR="$PROJECT_ROOT/pkg/rpc/proto"

# 确保输出目录存在
mkdir -p "$OUT_DIR/price" "$OUT_DIR/rsi"

# 生成Go代码
protoc --go_out="$OUT_DIR" --go-grpc_out="$OUT_DIR" \
    -I"$PROTO_DIR" \
    "$PROTO_DIR/price/price.proto" \
    "$PROTO_DIR/rsi/rsi.proto"

echo "Proto生成完成"
