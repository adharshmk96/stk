VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
MAJOR := $(shell echo $(VERSION) | cut -d . -f 1 | sed 's/v//')
MINOR := $(shell echo $(VERSION) | cut -d . -f 2)
PATCH := $(shell echo $(VERSION) | cut -d . -f 3)
NEW_PATCH := $(shell echo $$(($(PATCH) + 1)))
NEW_MINOR := $(shell echo $$(($(MINOR) + 1)))
NEW_MAJOR := $(shell echo $$(($(MAJOR) + 1)))
NEW_TAG_PATCH := v$(MAJOR).$(MINOR).$(NEW_PATCH)
NEW_TAG_MINOR := v$(MAJOR).$(NEW_MINOR).0
NEW_TAG_MAJOR := v$(NEW_MAJOR).0.0

current_branch := $(shell git branch --show-current)

.PHONY: patch minor major build test publish

##########################
### Manage Commands
##########################

patch:
	$(eval NEW_TAG := $(NEW_TAG_PATCH))
	$(call tag)

minor:
	$(eval NEW_TAG := $(NEW_TAG_MINOR))
	$(call tag)

major:
	$(eval NEW_TAG := $(NEW_TAG_MAJOR))
	$(call tag)

publish:
	@git push origin $(VERSION)

build:
	@go build .	

test:
	@go test -v ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

clean-branch:
	@git branch | grep -v "main" | xargs git branch -D


set_upstream:
	@current_branch=$(shell git branch --show-current); git branch --set-upstream-to=origin/$$current_branch $$current_branch



##########################
### Helpers
##########################
define tag
	@echo "current version is $(VERSION)"
    $(eval EXISTING_TAG := $(shell git tag -l $(NEW_TAG) 2>/dev/null))
    @if [ "$(EXISTING_TAG)" = "$(NEW_TAG)" ]; then \
        echo "Tag $(NEW_TAG) already exists. reapplying the tag."; \
        git tag -d $(NEW_TAG); \
    fi
    @git tag $(NEW_TAG)
    @echo "created new version $(NEW_TAG)."
endef

