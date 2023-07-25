##########################
### Version Commands
##########################

patch:
	$(eval NEW_TAG := $(shell git semver patch --dryrun))
	$(call update_file)
	@git semver patch

minor:
	$(eval NEW_TAG := $(shell git semver minor --dryrun))
	$(call update_file)
	@git semver minor

major:
	$(eval NEW_TAG := $(shell git semver major --dryrun))
	$(call update_file)
	@git semver major

publish:
	@git push origin $(shell git semver get)


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
### Helpers
##########################

define update_file
    @echo "updating files to version $(NEW_TAG)"
    @sed -i.bak "s/var version = \"[^\"]*\"/var version = \"$(NEW_TAG)\"/g" ./cmd/root.go
    @rm cmd/root.go.bak
    @git add cmd/root.go
    @git commit -m "bump version to $(NEW_TAG)" > /dev/null
endef

##########################
### configuration
##########################

init:
	@git config core.hooksPath .githooks
