build:
	go build -o ./bin/projectx


run: build
	./bin/projectx

tests:
	go test -v ./...

test:
	go test ./...

all: run
	echo "Hello"

test_coverage:
	@mkdir -p coverage
	@go test -coverprofile=coverage/coverage.out ./...

cover: test_coverage
	@go tool cover -func=coverage/coverage.out > coverage/coverage-summary.txt

clean_sumary: print_coverage	
	@rm ./coverage/coverage-summary.txt 

coverage: clean_sumary
	@rm ./coverage/coverage.out

print_coverage: cover
	@coverage=$$(go tool cover -func=coverage/coverage.out | grep total | awk '{print substr($$3, 1, length($$3)-1)}'); \
	echo "Test coverage: $$coverage%" \

