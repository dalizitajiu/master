echo "start building..."
set GOPATH=%GOPATH%;%cd%
go build %cd%/src/cache
go build %cd%/src/lib
go build %cd%/src/route
echo "end building"
go run %cd%/src/launch/launch.go
echo "start running"