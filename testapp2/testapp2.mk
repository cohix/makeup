

build:
	go build -o ${BIN_DEST}

run:
	${BIN_DEST}

test:
	go test -v ./...

env:
	echo "SOMETHING_IMPORTANT=another important value"

clean:
	rm ${BIN_DEST}