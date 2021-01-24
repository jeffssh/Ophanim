@echo off

echo building gRPC code
cd message && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative message.proto && cd ..
echo building all go bins
go install ./...
echo done!
