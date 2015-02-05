default: install

install:
	go generate
	go install .
