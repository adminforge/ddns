FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN ls 

RUN CGO_ENABLED=0 GOOS=linux go build -o ddns

FROM gcr.io/distroless/static-debian12 AS run

WORKDIR /

COPY --from=build --chown=65532:65532 /app/ddns /app/ddns

USER nonroot:nonroot

ENV GIN_MODE=release
ENV DDNS_EXPIRATION_DAYS=10

CMD ["/app/ddns"] 
