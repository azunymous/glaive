FROM golang:1.15
WORKDIR /igiari-glv
COPY . /igiari-glv/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ./cmd/glaive

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /igiari-glv/
COPY configfiles/production/config.yml ./config.yml
COPY --from=0 /igiari-glv/glaive .
CMD ["./glaive"]
