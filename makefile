publish:
	@git push && semver push


##########################
### Build Commands
##########################

BINARY_NAME=stk

build:
	@go build -o ./out/$(BINARY_NAME) -v

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
	@rm -f stk.db

deps:
	@go mod download

tidy:
	@go mod tidy

vet:
	@go vet

clean-branch:
	@git branch | egrep -v "(^\*|main|master)" | xargs git branch -D


##########################
### configuration
##########################

init:
	@git config core.hooksPath .githooks
	@go install github.com/adharshmk96/semver
	@go install github.com/vektra/mockery/v2@v2.35.4


##########################

template:
	@echo "Running project template generator script"
	@cd ./__template__ && python3 generate.py
