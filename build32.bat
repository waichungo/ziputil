set GOARCH=386
set CGO_ENABLED=1
go.exe build -buildvcs=false -o ziputil.exe -ldflags "-H=windowsgui  -linkmode external -extldflags '-static' -s -w"
set GOARCH=amd64
pause