.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/ping service/entrypoint/functions/ping/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/sqs_worker service/entrypoint/functions/sqs_worker/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/s3_worker service/entrypoint/functions/s3_worker/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
