FROM golang:1.23-bookworm AS build

RUN mkdir /app
WORKDIR /app
COPY . .

RUN go test ./...
RUN go build ./cmd/spectacles

###

FROM debian:bookworm-slim
COPY --from=build /app/spectacles /app/spectacles
ENTRYPOINT ["/app/spectacles"]


