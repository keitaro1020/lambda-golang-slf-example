#!/bin/sh

mockgen -source service/application/app.go -destination service/application/mocks/app.go -package mocks
mockgen -source service/domain/cat.go -destination service/domain/mocks/cat.go -package mocks
mockgen -source service/domain/repository.go -destination service/domain/mocks/repository.go -package mocks
mockgen -source service/domain/s3.go -destination service/domain/mocks/s3.go -package mocks
