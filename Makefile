.PHONY: all test clean
all: build fmt vet lint test

APP=goauction
ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")
UNIT_TEST_PACKAGES=$(shell  go list ./... | grep -v "vendor")

APP_EXECUTABLE="./out/$(APP)_server"

prepare:
	go get -u github.com/dgrijalva/jwt-go
	go get -u github.com/gin-gonic/gin
	go get -u github.com/iancoleman/strcase
	go get -u github.com/jinzhu/gorm
	go get -u syreclabs.com/go/faker
	go get -u github.com/golang/dep/cmd/dep
	go get -u golang.org/x/lint/golint
	go get -u github.com/axw/gocov/gocov
	go get -u gopkg.in/matm/v1/gocov-html
	go get -u github.com/swaggo/swag/cmd/swag
	go get -u github.com/rubenv/sql-migrate/...

migrate-init:
	cat dbconfig.yml.example > dbconfig.yml

migrate-new:
	sql-migrate new -env="development" $(name)

db-setup: db-create db-migrate

db-create:
	createdb -O$(DB_USER) -Eutf8 $(DB_NAME)

db-migrate:
	sql-migrate up -env="$(APP_ENV)"

db-drop:
	dropdb --if-exists -U$(DB_USER) $(DB_NAME)

db-reset: db-drop db-create db-migrate

testdb-migrate:
	sql-migrate up -env="test"

testdb-create: testdb-drop
	createdb -O$(DB_USER) -Eutf8 $(DB_NAME_TEST)

testdb-drop:
	dropdb --if-exists -U$(DB_USER) $(DB_NAME_TEST)

testdb-reset: testdb-drop testdb-create testdb-migrate

build-deps:
	dep ensure -v

update-deps:
	dep ensure

compile:
	mkdir -p out/
	go build -race -ldflags "-extldflags '-static'" $(APP_EXECUTABLE)

build: build-deps compile

install:
	go install ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	@for p in $(UNIT_TEST_PACKAGES); do \
		echo "==> Linting $$p"; \
		golint $$p | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } \
	done

test:
	go test -v ./tests -p=1

run-dev: 
	go run ${APP}.go

run:
	GIN_MODE=release go run ${APP}.go

api-docs:
	swag init -g goauction.go

copy-config:
	cp .env.example .env

clean:
	rm -rf ./out/ ./docs

