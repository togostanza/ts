default: install

install:
	npm run build
	go generate ./...
	go install .
