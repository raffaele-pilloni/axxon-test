include makefiles/mocks.mk

SERVICE_NAME=golang-http

start:
	@docker-compose up --build -d --remove-orphans --force-recreate

stop:
	@docker-compose stop

destroy:
	@docker-compose down --volumes --remove-orphans

restart: destroy start

ps:
	@docker-compose ps

logs:
	@docker-compose logs -f

service-logs:
	@docker-compose logs -f ${SERVICE_NAME}

terminal:
	@docker-compose exec ${SERVICE_NAME} sh

lint:
	@docker-compose exec golang-http golangci-lint run ./...

tests:
	@docker-compose exec golang-http sh -c "go run github.com/onsi/ginkgo/v2/ginkgo -r --keep-going --race"
