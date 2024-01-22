# Builder
FROM golang:alpine AS builder
LABEL stage=gobuilder
# ENVS
ENV CGO_ENABLED 0
ENV GOOS linux
# Copy
COPY . /src
WORKDIR /src
RUN go mod download && go build -o /build/theName ./cmd/app/main.go

# Productiobn
FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/theName /app/theName
CMD ["./theName"]