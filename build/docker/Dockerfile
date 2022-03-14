FROM golang:1.17.4-alpine AS build
WORKDIR /workspace

ARG app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o tmp/app cmd/${app}/*.go

FROM alpine:3.15.0
RUN apk --no-cache add ca-certificates curl
WORKDIR /workspace

COPY --from=build /workspace/tmp/app .
CMD ["./app"]
