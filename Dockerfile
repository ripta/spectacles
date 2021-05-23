FROM golang:1.16-stretch AS build
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


