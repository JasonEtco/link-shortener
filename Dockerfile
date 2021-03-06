FROM golang:1.11.2 as builder
WORKDIR /app
COPY go.mod go.sum *.go /app/
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o link-shortener .

FROM scratch
COPY --from=builder /app/link-shortener /
CMD ["/link-shortener"]
