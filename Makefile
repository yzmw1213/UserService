# Generete Code by Protocol Buffer
generate:
	sh generate.sh

# test go
test:
	docker-compose exec user_api go test -v ./grpc; \
	docker-compose exec user_api go test -v ./usecase/interactor

# golint
lint:
	docker-compose exec user_api golint ./...
