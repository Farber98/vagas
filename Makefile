tests:
	go test -cover -race -v -count=1 ./...

database:
	migrate -path internal/migrations -database "mysql://juan:juantuc98@tcp(localhost:3306)/pagarme_test?multiStatements=true" -verbose up