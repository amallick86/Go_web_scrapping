FROM golang:1.17.6 as builder
COPY go.mod go.sum /go/src/Go_web_scrapping/
WORKDIR /go/src/Go_web_scrapping
RUN go mod download
COPY . /go/src/Go_web_scrapping
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/webscrapping Go_web_scrapping

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/Go_web_scrapping/build/webscrapping /usr/bin/webscrapping
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/webscrapping"]
