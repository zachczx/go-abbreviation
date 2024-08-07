FROM golang:1.22

ENV GO111MODULE=on
# Set destination for COPY
WORKDIR /app

# Copy the directories with the local imports
COPY ./go.mod ./go.sum ./
COPY ./templates ./templates
COPY ./core ./core
COPY ./models ./models
COPY ./static ./static

# Download Go modules
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@$(go list -m -f '{{ .Version }}' github.com/a-h/templ)

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./

# Env
# https://github.com/coollabsio/coolify/issues/1918
ARG LISTEN_ADDR
ENV LISTEN_ADDR $LISTEN_ADDR 

# Build
RUN templ generate && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/go-abbreviation

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 5173

# Run
CMD ["/app/go-abbreviation"]