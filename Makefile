default: install

install:
	go-bindata -nocompress data/...
	go install .
