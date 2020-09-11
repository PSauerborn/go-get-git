FROM golang:latest as build

WORKDIR /app/server

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api/api.go

FROM alpine:latest as server

WORKDIR /app/server

COPY --from=build /app/server ./

EXPOSE 10071

CMD [ "./api" ]