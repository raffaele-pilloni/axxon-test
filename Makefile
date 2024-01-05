include makefiles/mocks.mk

SERVICE_NAME=golang-http

start:
	@docker-compose up --build -d --remove-orphans --force-recreate

stop:
	@docker-compose stop

destroy:
	@docker-compose down --volumes

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
	@docker-compose exec ${SERVICE_NAME} golangci-lint run ./...

test:
	@docker-compose exec -w /go/src/gitlab.facile.it/mutui/uxie/ ${SERVICE_NAME} sh -c "ENV=test go test ./tests/..."
