test:
	go test

test-coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

build:
	env GOOS=darwin GOARCH=amd64 go build -o ./bin/planespotter_darwin_amd64
	env GOOS=linux GOARCH=amd64 go build -o ./bin/planespotter_linux_amd64
	env GOOS=windows GOARCH=amd64 go build -o ./bin/planespotter_win_amd64.exe