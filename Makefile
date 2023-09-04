up:
	docker-compose -f docker-compose.yaml up

rm:
	docker-compose stop \
	&& docker-compose rm db \
	&& docker-compose rm app \
	&& rmdir pg-data /s /q

test:
	cd internal/controller/http/v1 \
	&& go test -v \
	&& cd ../../../../

swag:
	swag init -g cmd/app/main.go --parseInternal --parseDependency
