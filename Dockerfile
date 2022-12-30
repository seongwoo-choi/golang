FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux 

WORKDIR /build

COPY . .

RUN go mod download

RUN go build server.go

# RUN go build -o server .

WORKDIR /dist

RUN cp /build/server .

FROM scratch

COPY --from=builder /dist/server .

EXPOSE 3000

ENTRYPOINT ["/server"]