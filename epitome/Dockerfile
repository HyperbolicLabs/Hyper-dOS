###
# 1) build the binary
###
FROM golang:1.22 as builder

WORKDIR /app

COPY . .

RUN echo "ignoring these files: " && cat .dockerignore

RUN go mod download

COPY *.go ./

# CGO_ENABLED=0 allows us to create a statically-linked binary that can run without any external dependencies.
RUN CGO_ENABLED=0 go build -o /epitome


###
# 2) keep the binary, install requirements, discard the rest.
###

FROM alpine:latest

COPY --from=builder /epitome /epitome

ENTRYPOINT ["/epitome"]