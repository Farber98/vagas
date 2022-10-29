tests:
	go test -cover -race -v -count=1 ./...

database:
	migrate -path internal/migrations -database "mysql://juan:juan@tcp(mysql)/pagarme_test?multiStatements=true" -verbose up

.PHONY: tests database