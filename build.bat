@ECHO off

echo "--> Building Windows binary..."
cd src
go env -w GO111MODULE=on
go build -o ../bin/win/4bid-n.exe
cd ..