up:
	docker-compose -f docker-compose.yaml up

rm:
	docker-compose stop \
	&& docker-compose rm db \
	&& rmdir pg-data /s /q

swag:
	swag init -g cmd/app/main.go --parseInternal --parseDependency