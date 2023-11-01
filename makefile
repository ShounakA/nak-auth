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
