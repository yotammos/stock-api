set -xe

go get github.com/yotammos/stock-api

go build -o bin/application application.go