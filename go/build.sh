rm -f go.mod
rm -f go.sum
rm -f main

go mod init main
go mod tidy
go build main
