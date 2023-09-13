test:
	go test ./...

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

build-windows:
	fyne package -os windows --src src --executable ../bin/planespotter_win.exe

build-linux:
	fyne package -os linux --src src --executable ../bin/planespotter_linux

build-mac:
	fyne package -os darwin --src src --executable ../bin/planespotter_mac