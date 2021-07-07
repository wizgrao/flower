FROM golang:1.15
WORKDIR /app/server
COPY . .
RUN go build cmd/run.go
RUN ls

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /app/server/run .
COPY index.html index.html 
COPY flower.js flower.js
RUN apk add --no-cache libc6-compat
CMD ["/root/run"]

