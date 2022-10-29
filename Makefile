tests:
	go test -cover -race -v -count=1 ./...

ddl: 
	migrate -path internal/scripts/ddl.sql -database "mysql://juan:juan@tcp(localhost:3306)/pagarme_test" -verbose up