@echo off

echo building gRPC code
:: Begin building gRPC code
cd message
:: golang client and server 
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative message.proto
:: cpp client
protoc --cpp_out=cpp message.proto

echo building all go bins
go install ./...
echo done!
