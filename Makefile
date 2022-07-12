run:
	echo "running.."
	go run ./cmd/shortener

build:
	echo "building.."
	go build -a -o ./cmd/shortener ./cmd/shortener

test:
	echo "testing.."
	go test -v -cover ./...

.DEFAULT_GOAL := run