#build stage
FROM golang:1.17-alpine3.14 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -fsSL \
    https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
    sh 

#Run stage
FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /usr/local/bin/goose ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration


EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]