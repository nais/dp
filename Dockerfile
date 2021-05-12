FROM golang:1.13-alpine as builder
RUN apk add --no-cache git make
ENV GOOS=linux
ENV CGO_ENABLED=0
ENV GO111MODULE=on
COPY . /src
WORKDIR /src
RUN rm -f go.sum
RUN go get ./...
RUN make test
RUN make linux-build

FROM alpine:3
WORKDIR /app
COPY --from=builder /src/dp /app/dp
CMD ["/app/dp"]
