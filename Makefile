# Makefile

BINARY_NAME=go-abbreviation

run: build 
	@./bin/main.exe

build:
	go mod tidy && \
   	templ generate && \
	go generate && \
	go build -ldflags="-w -s" -o ./bin/main.exe

dev/templ:
	templ generate --watch --proxy="http://localhost:5173" --open-browser=false -v

dev/tailwind:
	npx tailwindcss -i ./static/css/input.css -o static/css/styles.css --minify --watch

dev/air:
	air -c .air.toml

dev/prettier:
	npx prettier . --write

dev: 
	make -j3 dev/templ dev/tailwind dev/air 
	
#dev/prettier