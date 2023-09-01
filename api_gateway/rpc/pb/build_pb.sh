#!/bin/bash

# 指定输出目录
output_dir="./"

# 列出所有的 .proto 文件
proto_files=(
  "user.proto"
  "video.proto"
  "relation.proto"
  "favorite.proto"
  "message.proto"
  "comment.proto"
)

for proto_file in "${proto_files[@]}"
do
  echo "Compiling $proto_file..."

  protoc --go_out=. --go_opt=paths=source_relative \
         --go-grpc_out=. --go-grpc_opt=paths=source_relative "$proto_file"

  echo "Compilation of $proto_file completed."
done

echo "All .proto files compiled successfully."
