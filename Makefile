default: build

build:
	go-bindata -nocompress data/...
	go build -o togostanza .
