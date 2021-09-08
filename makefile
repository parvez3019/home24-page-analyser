IMAGE_TAG=home24-page-analyser

.PHONY: test
test:
	go test -cover ./...

.PHONY: build
build:
	go build -a -installsuffix cgo -o build/server cmd/server/main.go

.PHONY: docker-build
docker-build:
	docker build -t home24-page-analyser .

.PHONY: run
docker-run: docker-build
	docker run \
 		-p 8000:8000 \
		-e PORT=8000 \
		--rm -it $(IMAGE_TAG)
