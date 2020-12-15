.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/ping cmd/functions/ping/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/sqs_worker cmd/functions/sqs_worker/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/s3_worker cmd/functions/s3_worker/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

mockgen:
	chmod u+x mockgen.sh
	./mockgen.sh
