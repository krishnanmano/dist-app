GOLANGCI_LINT_VERSION=v1.44
DB_URL=postgresql://postgres:pgsecret@localhost:55432/catalog_service?sslmode=disable

.PHONY: list
list:
	@sh -c "$(MAKE) -p no_targets__ 2>/dev/null | \
        awk -F':' '/^[a-zA-Z0-9][^\$$#\/\\t=]*:([^=]|$$)/ {split(\$$1,A,/ /);for(i in A)print A[i]}' | \
        grep -v Makefile | \
        grep -v '%' | \
        grep -v '__\$$' | \
        sort -u"

.PHONY: go-lint-docker
go-lint-docker:
	docker run --rm \
	-v ${PWD}:/tmp/lint \
	-w /tmp/lint \
	golangci/golangci-lint:$(GOLANGCI_LINT_VERSION)-alpine golangci-lint run

.PHONY: shellcheck
shellcheck:
	docker run --rm \
	-v ${PWD}:/tmp/lint \
	-w /tmp/lint \
	koalaman/shellcheck:latest shellcheck deploy.sh

.PHONY: lint
lint:
	golangci-lint run

.PHONY: mac_amd64_build
mac_amd64_build:
	env GOOS=darwin GOARCH=amd64 go build -o distapp .

.PHONY: linux_build
linux_build:
	env GOOS=linux GOARCH=amd64 go build -o distapp .

.PHONY: docker_build
docker_build:
	docker build -t distapp:0.0.1 .

.PHONY: distapp
distapp:


.PHONY: delete_docker
delete_docker: