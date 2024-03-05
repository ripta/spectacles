FROM golang:1.21-bullseye AS build
ENV GO111MODULE=on

RUN mkdir /app
WORKDIR /app
COPY . .

RUN go test ./...
RUN go build ./cmd/spectacles

###

FROM debian:stretch-slim
COPY --from=build /app/spectacles /app/spectacles
ENTRYPOINT ["/app/spectacles"]


