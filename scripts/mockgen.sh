#!/bin/sh

mockgen -source ../pkg/application/app.go -destination ../pkg/application/mocks/app.go -package mocks
mockgen -source ../pkg/domain/cat.go -destination ../pkg/domain/mocks/cat.go -package mocks
mockgen -source ../pkg/domain/repository.go -destination ../pkg/domain/mocks/repository.go -package mocks
mockgen -source ../pkg/domain/s3.go -destination ../pkg/domain/mocks/s3.go -package mocks
