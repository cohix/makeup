
build:
	go build -o ${BIN_DEST}

run:
	${BIN_DEST}

test:
	go test -v ./...

env:
	echo "CONFIG_KEY=some_config_val"

clean:
	rm ${BIN_DEST}
