test:
	go test ./...

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

build:
	fyne package -os windows --executable ./bin/planespotter_win.exe