mock_dal:
	@docker-compose exec $(SERVICE_NAME) mockery --all --keeptree --case snake --dir=./internal/dal

mock_client:
	@docker-compose exec $(SERVICE_NAME) mockery --all --keeptree --case snake --dir=./internal/client

mock_service:
	@docker-compose exec $(SERVICE_NAME) mockery --all --keeptree --case snake --dir=./internal/service

mock_repository:
	@docker-compose exec $(SERVICE_NAME) mockery --all --keeptree --case snake --dir=./internal/repository

mock: mock_dal mock_client mock_service mock_repository
