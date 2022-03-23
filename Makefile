.PHONY:install
install:
	@go install github.com/golang/mock/mockgen@v1.6.0
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.2

.PHONY:mock
mock:
	@mockgen --package=repository --source=pkg/persistence/repository/account.go --destination=pkg/persistence/repository/account_mock.go Accounts
	@mockgen --package=repository --source=pkg/persistence/repository/operation.go --destination=pkg/persistence/repository/operation_mock.go Operations
	@mockgen --package=repository --source=pkg/persistence/repository/transaction.go --destination=pkg/persistence/repository/transaction_mock.go Transactions
	@mockgen --package=service --source=pkg/service/transaction.go --destination=pkg/service/transaction_mock.go Transactions
	@mockgen --package=service --source=pkg/service/account.go --destination=pkg/service/account_mock.go Accounts
	@mockgen --package=common --source=pkg/common/log.go --destination=pkg/common/log_mock.go Logger

.PHONY:lint
lint:
	@golangci-lint -v run ./...

.PHONY:test
test:
	@go clean -testcache
	@go test -v ./... | grep -v "_mock.go"

.PHONY:coverage
coverage:
	@go clean -testcache
	@go test -coverprofile=tmp/coverage.out ./...
	@go tool cover -html=tmp/coverage.out

.PHONY:run
run:
	@docker-compose up --build

.PHONY:seeds
seeds:
	@exec scripts/seeds.sh