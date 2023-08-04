BIN=golua
BIN_SRC=main.go
DIST_DIR=.dist
COVERAGE_FILE=.coverage
LOG_FILE=*.log*



##################################
#######       Setup       ########
##################################
.PHONY: ensure install_lint

tidy:
	@go mod tidy

download: tidy
	@go mod download

install_lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3


##################################
#######        Tool       ########
##################################
.PHONY: fmt lint build clean api

fmt:
	@golangci-lint run --fix

lint:
	@golangci-lint run ./...

build:
	@go build -o ${DIST_DIR}/${BIN} ${BIN_SRC}

clean:
	@git clean -fdx ${DIST_DIR} ${COVERAGE_FILE} ${LOG_FILE}



##################################
#######       test        ########
##################################
TEST_FLAGS = -v -race -failfast -covermode=atomic
MINIMAL_COVERAGE = 40

.PHONY: test coverage coverage_html

test:
	@CGO_ENABLED=1 ROOT=${PWD} ENV=testing go test ${TEST_FLAGS} -coverprofile=${COVERAGE_FILE} -cover -timeout=40s `go list ./...`

coverage: test
	@go tool cover -func=${COVERAGE_FILE}
	@COVERAGE=$$(go tool cover -func=${COVERAGE_FILE} | grep total | awk '{print $$3}' | sed 's/\%//g'); \
	echo "Current coverage is $${COVERAGE}%, minimal is ${MINIMAL_COVERAGE}."; \
	awk "BEGIN {exit ($${COVERAGE} < ${MINIMAL_COVERAGE})}"

coverage_html: test
	@go tool cover -html ${COVERAGE_FILE}


##################################
#######      CI/CD        ########
##################################

.PHONY: ci_init ci_test 

ci_init: download install_lint 

ci_test: lint coverage


