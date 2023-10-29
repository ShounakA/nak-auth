build:
	make buildLoginPage; \
	cd src; \
	go build .;
run:
	cd src; \
	./bin/air go run .;
install:
	cd src; \
	go get .
buildLoginPage:
	rm -rf src/static; \
	mkdir src/static; \
	cd ../frontend; \
	yarn build; \
	cp -r dist/* ../backend/src/static;
dbDevProxy:
	pscale connect nak-data staging
