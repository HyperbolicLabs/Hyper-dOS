###
# 1) build the binary
###
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.24 AS builder

ARG TARGETARCH
ENV GOOS=linux GOARCH=$TARGETARCH

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

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:latest
# FROM --platform=${TARGETPLATFORM:-linux/amd64} gcr.io/distroless/static:nonroot

COPY --from=builder /epitome /epitome
COPY ./application.yaml /application.yaml
# USER nonroot:nonroot

ENTRYPOINT ["/epitome"]
