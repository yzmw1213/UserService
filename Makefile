# Generete Code by Protocol Buffer
generate:
	sh generate.sh

# test go
test:
	docker-compose exec blog_api go test -v ./grpc

# golint
lint:
	docker-compose exec blog_api golint ./...
