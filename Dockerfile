FROM golang:1.22.5-bookworm AS first
ENV GO111MODULE=on
WORKDIR /app
COPY ./go.mod ./go.sum tailwind.config.js package.json package-lock.json ./abbreviations.db ./
COPY ./templates ./templates
COPY ./core ./core
COPY ./models ./models
COPY ./metaphone3 ./metaphone3
COPY ./static ./static
# Download Go modules
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@$(go list -m -f '{{ .Version }}' github.com/a-h/templ)
# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./
# Env --- https://github.com/coollabsio/coolify/issues/1918
ARG LISTEN_ADDR
ENV LISTEN_ADDR $LISTEN_ADDR 
# Build
RUN templ generate && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/go-abbreviation

#####################

FROM node:22 AS second
WORKDIR /app
COPY --from=first /app/tailwind.config.js /app/abbreviations.db /app/package.json /app/go-abbreviation /app/package-lock.json /app/
COPY --from=first /app/templates /app/templates
COPY --from=first /app/static /app/static
# COPY package*.json ./
RUN npm install
RUN npx tailwindcss -i /app/static/css/input.css -o /app/static/css/styles.css --minify
RUN npx brotli-cli compress --glob /app/static/css/styles.css /app/static/js/htmx.min.js /app/static/data.json

#####################

FROM scratch
#FROM alpine:3.20
WORKDIR /app
COPY --from=second /app/go-abbreviation ./go-abbreviation
COPY --from=second /app/static ./static
COPY --from=second /app/abbreviations.db ./abbreviations.db
EXPOSE 5173

# Run
CMD ["/app/go-abbreviation"]

