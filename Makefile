default: install

install:
	bower install
	go generate ./...
	go install .
