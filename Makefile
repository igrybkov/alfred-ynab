# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./
BINARY_NAME := workflow
VERSION := $(shell /usr/libexec/PlistBuddy -c 'Print :version' info.plist)
WORKFLOW_FILE_NAME := $(shell /usr/libexec/PlistBuddy -c 'Print :name' info.plist)-$(VERSION).alfredworkflow

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# .refresh-version: refresh version variable
.PHONY: .refresh-version
.refresh-version:
	$(eval VERSION := $(shell /usr/libexec/PlistBuddy -c 'Print :version' info.plist))
	$(eval WORKFLOW_FILE_NAME := $(shell /usr/libexec/PlistBuddy -c 'Print :name' info.plist)-$(VERSION).alfredworkflow)

# .check-dependencies: check if dependencies are installed
.PHONY: .refresh-version
.check-dependencies:
	$(if $(shell command -v go), , $(error "go is not installed. Please install it from https://golang.org/dl/"))
	$(if $(shell command -v gh), , $(error "gh is not installed. Please install it from https://cli.github.com/"))
	$(if $(shell command -v lipo), , $(error "lipo is not installed. You may build on an unsupported platform."))
	$(if $(shell command -v zip), , $(error "zip is not installed. Please install it before continuing."))

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

## build: build the application
.PHONY: build
build:
	@GOOS=darwin GOARCH=amd64 go build -o ${BINARY_NAME}_amd64 ${MAIN_PACKAGE_PATH}
	@GOOS=darwin GOARCH=arm64 go build -o ${BINARY_NAME}_arm64 ${MAIN_PACKAGE_PATH}
	@lipo -create -output ${BINARY_NAME} ${BINARY_NAME}_amd64 ${BINARY_NAME}_arm64
	@rm ${BINARY_NAME}_amd64 ${BINARY_NAME}_arm64

## workflow: build the alfred workflow
.PHONY: workflow
workflow: set-version update-readme clean build
	@zip -r "${WORKFLOW_FILE_NAME}" ${BINARY_NAME}* README.md LICENSE info.plist icons icon.png images run.sh
	@echo "Workflow created: ${WORKFLOW_FILE_NAME}"


## clean: remove build artifacts
.PHONY: clean
clean:
	@rm -f ${BINARY_NAME} ${BINARY_NAME}_amd64 ${BINARY_NAME}_arm64 *.alfredworkflow

# ==================================================================================== #
# RELEASE
# ==================================================================================== #

## set-version: set workflow version to plist
.PHONY: set-version
set-version:
	@read -p "Enter version [$(VERSION)]: " VERSION \
		&& VERSION="$${VERSION:-"$(VERSION)"}" \
		&& /usr/libexec/PlistBuddy -c "Set :version $$VERSION" info.plist && echo "Version updated"

## update-readme: set workflow readme to plist
.PHONY: update-readme
update-readme:
	@/usr/libexec/PlistBuddy -c "Import :readme README.md" info.plist && echo "Readme updated"

## tag-release: create release tag from current version
.SILENT: tag-release
tag-release:
	git tag v$(VERSION)

## tag-push: push current release tag
.SILENT: tag-push
tag-push:
	git push origin v$(VERSION)

# .check-git-status: check if git status is clean
.PHONY: .check-git-status
.check-git-status:
	$(eval GIT_STATUS := $(shell git status --porcelain))
	$(if $(GIT_STATUS), $(error "Working directory is not clean. Please commit or stash your changes."))

# .check-version: check if version already exists
.PHONY: .check-version
.check-version:
	$(if $(shell git tag -l v$(VERSION)), $(error "Version $(VERSION) already exists. Please update the version in info.plist."))

## release: create a new release
.PHONY: release
release: .check-dependencies .check-git-status set-version .refresh-version .check-version workflow tag-release tag-push
	@gh release create --draft --generate-notes --latest --verify-tag "v$(VERSION)" *.alfredworkflow
	@echo "Release v$(VERSION) created"
