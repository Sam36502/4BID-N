@ECHO off

echo "--> Building Windows binary..."
go build -o bin/win/4bod.exe src\main.go src\screen.go src\fbod.go src\options.go src\input.go