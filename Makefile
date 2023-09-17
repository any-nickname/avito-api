up:
	docker-compose -f docker-compose.yaml up

stop:
	docker-compose stop

down:
	docker-compose down -v \
	&& docker rmi avito-api-app \
	&& rmdir pg-data /s /q

test:
	cd internal/controller/http/v1 \
	&& go test -v \
	&& cd ../../../../

swag:
	swag init -g cmd/app/main.go --parseInternal --parseDependency
