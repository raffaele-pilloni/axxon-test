mock_dal:
	@docker-compose exec $(SERVICE_NAME) mockery --all --keeptree --case snake --output mock --outpkg mock --dir ./internal/db

mock_client:
	@docker-compose exec $(SERVICE_NAME) mockery --all --keeptree --case snake --output mock --outpkg mock --dir ./internal/client

mock_service:
	@docker-compose exec $(SERVICE_NAME) mockery --all --keeptree --case snake --output mock --outpkg mock --dir ./internal/service

mock_repository:
	@docker-compose exec $(SERVICE_NAME) mockery --all --keeptree --case snake --output mock --outpkg mock --dir ./internal/repository

mock_error:
	@docker-compose exec $(SERVICE_NAME) mockery --all --keeptree --case snake --output mock --outpkg mock --dir ./internal/error

mock: mock_dal mock_client mock_service mock_repository mock_error
