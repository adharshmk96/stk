publish:
	@git push && semver push


##########################
### Build Commands
##########################

BINARY_NAME=app

build:
	@go build -o ./out/$(BINARY_NAME) -v

run: 
	@go run . serve -p 8080

test:
	@go test ./... -coverprofile=coverage.out

coverage:
	@go test -v ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

testci:
	@go test ./... -coverprofile=coverage.out

clean:
	@go clean
	@rm -f ./out/$(BINARY_NAME)
	@rm -f coverage.out
	@rm -rf .keys
	@rm -f auth_database.db

deps:
	@go mod download

tidy:
	@go mod tidy

lint:
	@golangci-lint run --enable-all

vet:
	@go vet

clean-branch:
	@git branch | egrep -v "(^\*|main|master)" | xargs git branch -D

	
##########################
### Setup Commands
##########################

init: deps keygen initgithooks mockgen
	@echo "Project initialized."

initci: deps keygen
	@echo "Project initialized for CI."

initgithooks:
	@git config core.hooksPath .githooks

mockgen:
	@rm -rf ./mocks
	@mockery --all	

install-tools:
	@go install github.com/adharshmk96/semver
	@go install github.com/vektra/mockery/v2@v2.35.4

