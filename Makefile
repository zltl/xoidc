
.PHONY: clean deps jet

.ONESHELL:


deps:
	# if jet command not found, install
	if ! [ -x "$$(command -v jet)" ]; then
		echo "Installing jet...";
		go install github.com/go-jet/jet/v2/cmd/jet@latest;
	fi

jet: deps
	echo "Generating jet db operation files..."
	jet -dsn=postgresql://postgres:123456@localhost:5432/xoidc?sslmode=disable -schema=public -path=./.gen

clean:
	rm -rf bin/*

