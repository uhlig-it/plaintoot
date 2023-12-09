FROM golang
WORKDIR /go/src/github.com/uhlig-it/plaintoot
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o plaintoot .

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/uhlig-it/plaintoot/plaintoot /usr/local/bin
CMD ["./plaintoot", "--help"]
