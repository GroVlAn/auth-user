FROM golang:1.25.1 as builder

ARG CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o user ./cmd/main.go 

FROM scratch
COPY --from=builder /app/user /user
COPY .env .
COPY ./configs ./configs
EXPOSE 8081

ENTRYPOINT ["/user"]