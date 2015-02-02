default: build

build:
	go-bindata data/...
	go build -o togostanza .
