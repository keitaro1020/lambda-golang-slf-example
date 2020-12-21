.PHONY: build clean deploy gomodgen

build: gomodgen graphqlgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/ping cmd/functions/ping/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/sqs_worker cmd/functions/sqs_worker/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/s3_worker cmd/functions/s3_worker/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get_cat cmd/functions/get_cat/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/graphql cmd/functions/graphql/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

mockgen:
	chmod u+x scripts/mockgen.sh
	cd  scripts/ && ./mockgen.sh

test:
	go test ./...

graphqlgen:
	cd scripts/graphql && gqlgen gen