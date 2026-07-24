# Dev Dockerfile — just runs the code with `go run`, no compiled binary.
# Good for learning/dev; slower startup and bigger image than a
# multi-stage build, but simplest to reason about and rebuilds fast
# once module caching is wired up later.
FROM golang:1.26-alpine

WORKDIR /app

# Copy dependency files first so Docker can cache this layer.
# As long as go.mod/go.sum don't change, "go mod download" is skipped
# on rebuilds — only your actual code changes trigger a re-download.
COPY go.mod go.sum ./
RUN go mod download

# Now copy the rest of your source code (cmd/, internal/, migrations/, web/).
COPY . .

EXPOSE 8080

# Runs your entrypoint directly from source every time the container starts.
# Change "./cmd/api" if your main package lives somewhere else.
CMD ["go", "run", "./cmd/api"]