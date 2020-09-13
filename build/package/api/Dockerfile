FROM golang:latest as build

WORKDIR /app/server

COPY ./api.go ./

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build ./api.go

FROM alpine:latest as server

WORKDIR /app/server

COPY --from=build /app/server ./

RUN chmod +x ./api

EXPOSE 10071

CMD [ "./api" ]