build:
	cd src; \
	go build .;
run:
	cd src; \
	./bin/air go run .;
install:
	cd src; \
	go get .
dbDevProxy:
	pscale connect nak-data staging
test:
	cd src; \
	go test ./... -v -coverprofile=coverage.out; \
	go tool cover -html=coverage.out -o coverage.html; \
	go tool cover -func=coverage.out;