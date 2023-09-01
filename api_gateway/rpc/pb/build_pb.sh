#!/bin/bash
###
 # @Autor: violet apricity ( Zhuangpx )
 # @Date: 2023-09-01 14:03:04
 # @LastEditors: violet apricity ( Zhuangpx )
 # @LastEditTime: 2023-09-01 14:56:33
 # @FilePath: \Road2TikTok\api_gateway\rpc\pb\build_pb.sh
 # @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
### 

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

  protoc --proto_path=. --go_out=. --go_opt=paths=source_relative \
         --go-grpc_out=. --go-grpc_opt=paths=source_relative "$proto_file"

  echo "Compilation of $proto_file completed."
done

echo "All .proto files compiled successfully."
