FROM golang:1.18.0-alpine as build

RUN apk update && apk upgrade && apk add --no-cache git

WORKDIR /app/student-api

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./out/api .

FROM alpine:latest

# Adds CA Certificates to the image
RUN apk add ca-certificates
COPY --from=build /app/student-api/out/api /app/student-api
COPY --from=build /app/student-api/.env /app/.env



WORKDIR "/app"

EXPOSE 9091

ENTRYPOINT ["./student-api"]
