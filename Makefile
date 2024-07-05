BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')
DOCKER := $(shell which docker)

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --exact-match 2>/dev/null)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

# Update the ldflags with the app, client & server names
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=mini \
	-X github.com/cosmos/cosmos-sdk/version.AppName=minid \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

###########
# Install #
###########

all: install

install:
	@echo "--> ensure dependencies have not been modified"
	# @go mod verify
	@echo "--> installing minid"
	@go install $(BUILD_FLAGS) -mod=readonly ./cmd/minid

init:
	./scripts/init.sh

###############################################################################
###                                Linting                                  ###
###############################################################################

golangci_version=v1.59.0

#? setup-pre-commit: Set pre-commit git hook
setup-pre-commit:
	@cp .git/hooks/pre-commit .git/hooks/pre-commit.bak 2>/dev/null || true
	@echo "Installing pre-commit hook..."
	@ln -sf ../../scripts/hooks/pre-commit.sh .git/hooks/pre-commit
	@echo "Pre-commit hook installed successfully"

#? lint-install: Install golangci-lint
lint-install:
	@echo "--> Installing golangci-lint $(golangci_version)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)

#? lint: Run golangci-lint
lint:
	@echo "--> Running linter"
	$(MAKE) lint-install
	@./scripts/go-lint-all.bash --timeout=15m

#? lint: Run golangci-lint and fix
lint-fix:
	@echo "--> Running linter"
	$(MAKE) lint-install
	@./scripts/go-lint-all.bash --fix

.PHONY: lint lint-fix

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.14.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

#? proto-all: Run make proto-format proto-lint proto-gen
proto-all: proto-format proto-lint proto-gen

#? proto-gen: Generate Protobuf files
proto-gen:
	@$(protoImage) sh ./x/crs/protocgen.sh

#? proto-swagger-gen: Generate Protobuf Swagger
proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@$(protoImage) sh ./scripts/protoc-swagger-gen.sh

#? proto-format: Format proto file
proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

#? proto-lint: Lint proto file
proto-lint:
	@$(protoImage) buf lint --error-format=json

#? proto-check-breaking: Check proto file is breaking
proto-check-breaking:
	@$(protoImage) buf breaking --against $(HTTPS_GIT)#branch=main

#? proto-update-deps: Update protobuf dependencies
proto-update-deps:
	@echo "Updating Protobuf dependencies"
	$(DOCKER) run --rm -v $(CURDIR)/proto:/workspace --workdir /workspace $(protoImageName) buf mod update

.PHONY: proto-all proto-gen proto-swagger-gen proto-format proto-lint proto-check-breaking proto-update-deps
