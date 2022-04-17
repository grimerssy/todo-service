run:
	go run ./cmd/main.go


start:
	brew services start postgresql

stop:
	brew services stop postgresql

create:
	createdb todo

drop:
	dropdb todo