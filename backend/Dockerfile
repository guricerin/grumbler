FROM golang:1.17 AS builder

WORKDIR /go/src/app
# copy the go src
COPY . .
RUN go build -o /go/bin/app ./cmd/app/main.go

FROM gcr.io/distroless/base-debian11 AS production

COPY --from=builder /go/bin/app /
CMD ["/app"]
