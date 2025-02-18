SRC=.
BIN=hmsv

build:
	go build $(SRC)

linux:
	GOOS=linux GOARCH=amd64 go build -o linux-amd64/$(BIN) $(SRC)

win:
	GOOS=windows GOARCH=amd64 go build -o windows-amd64/$(BIN).exe $(SRC)

run:
	go run $(SRC)

test:
	go test -run ''

lint:
	golangci-lint run
