
.PHONY: clean deps jet build

.ONESHELL:

build:
	echo "Building xoidc..."

	go build -o bin/xoidc ./cmd/xoidc_server

# if jet command not found, install
deps:
	if ! [ -x "$$(command -v jet)" ]; then
		echo "Installing jet...";
		go install github.com/go-jet/jet/v2/cmd/jet@latest;
	fi

jet: deps
	echo "Generating jet db operation files..."
	jet -dsn=postgresql://postgres:123456@localhost:5432/xoidc?sslmode=disable -schema=public -path=./gen

bksql:
	echo "Backup sql files..."
	export PGPASSWORD=postgres
	pg_dump -U postgres -d xoidc -f ./xoidc.sql --host localhost

clean:
	rm -rf bin/*

