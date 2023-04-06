set GOARCH=386
set CGO_ENABLED=1
go build -buildmode=c-shared -o ziputil32.dll -ldflags "-s -w" main.go
set GOARCH=amd64