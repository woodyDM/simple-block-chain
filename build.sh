#!/bin/sh

SHELL_FOLDER=$(cd "$(dirname "$0")";pwd)
UI=$SHELL_FOLDER/ui
echo "use working folder $SHELL_FOLDER"
cd $SHELL_FOLDER
rm -rf $SHELL_FOLDER/build
mkdir $SHELL_FOLDER/build
mkdir $SHELL_FOLDER/build/static
echo "start to build golang"
go build -o build/app cmd/server/main.go
cd $UI
echo `pwd`
echo "Start to compile ui"
npm install
npm run build
cd $SHELL_FOLDER
echo "当前目录 ：`pwd`,拷贝文件"
cp -rf ui/build/* build/static/
cp *.json build/
cp Dockerfile build
echo "Package complete "
