FROM golang:1.25 AS build

WORKDIR /app

COPY ./build/go.mod ./build/go.sum ./

RUN go mod download

COPY ./build/ ./

RUN ls 

RUN CGO_ENABLED=0 GOOS=linux go build -o ddns

FROM gcr.io/distroless/static-debian12 AS run

WORKDIR /

COPY --from=build --chown=65532:65532 /app/ddns /app/ddns
COPY ./build/static /app/static
USER nonroot:nonroot

ENV GIN_MODE=release
ENV DDNS_EXPIRATION_DAYS=90

CMD ["/app/ddns"] 
