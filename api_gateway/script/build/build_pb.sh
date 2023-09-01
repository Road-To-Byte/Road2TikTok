###
 # @Autor: violet apricity ( Zhuangpx )
 # @Date: 2023-09-01 14:00:46
 # @LastEditors: violet apricity ( Zhuangpx )
 # @LastEditTime: 2023-09-01 14:01:43
 # @FilePath: \Road2TikTok\api_gateway\script\build\build_pb.sh
 # @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
### 

#!/bin/bash

# 指定输出目录
output_dir="../"

# 列出所有的 .proto 文件
proto_files=(
  "pb/xx.proto"
  "pb/yy.proto"
  "pb/zz.proto"
  # 添加更多 .proto 文件
)

# 循环编译每个 .proto 文件
for proto_file in "${proto_files[@]}"
do
  echo "Compiling $proto_file..."
  
  # 使用 protoc 编译 .proto 文件
  protoc --go_out="$output_dir" --go_opt=paths=source_relative \
         --go-grpc_out="$output_dir" --go-grpc_opt=paths=source_relative "$proto_file"
  
  echo "Compilation of $proto_file completed."
done

echo "All .proto files compiled successfully."
