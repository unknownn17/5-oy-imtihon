run:
	go run cmd/main.go
table:
	migrate create -dir migrations -ext sql db
up:
	migrate -path migrations -database "postgres://postgres:2005@localhost/hotels?sslmode=disable" up
down:
	migrate -path migrations -database "postgres://postgres:2005@localhost/hotels?sslmode=disable" down
force:
	migrate -path migrations -database "postgres://postgres:2005@localhost/hotels?sslmode=disable" force 