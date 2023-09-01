#!/bin/bash

search_dir="./"

file_extension=".pb.go"

find "$search_dir" -type f -name "*$file_extension" -exec rm {} \;

echo "Files with '$file_extension' extension in '$search_dir' have been deleted."
