#!/bin/bash

pb_dir=$(pwd)/pb
cd proto
proto_dir=$(pwd)
go_path=$(echo $GOPATH)

file_names=$(ls $proto_dir)
for file_name in $file_names
do
    file_path=$proto_dir/$file_name
    protoc --proto_path="$go_path/src/" -I . --go_out=plugins=grpc:$pb_dir $file_name
    code=$(echo $?)
    if [ $code = 0 ]; then
        echo -e "编译文件: $file_path => [\033[31m成功\033[0m] "
    else
        echo -e "\t 编译文件: [$file_path] => [\033[34m失败\033[0m] "
        echo -e "\t <<<<<<<<<<<< 编译过程意外退出，已终止  <<<<<<<<<<<<"
        exit
    fi
done