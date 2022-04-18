run:
	go run ./cmd/main.go

test:
	go test ./...

tidy:
	go mod tidy

start:
	brew services start postgresql

stop:
	brew services stop postgresql

create:
	createdb todo

drop:
	dropdb todo

up:
	migrate -path ./schema -database "postgresql://USERNAME:PASSWORD@localhost:5432/todo?sslmode=disable" -verbose up

down:
	migrate -path ./schema -database "postgresql://USERNAME:PASSWORD@localhost:5432/todo?sslmode=disable" -verbose down