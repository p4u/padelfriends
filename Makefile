build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm ./cmd/padelfriends
	go build ./cmd/padelfriends

run: build
	./padelfriends
