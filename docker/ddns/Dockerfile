FROM node:22 AS sass

WORKDIR /theme

COPY theme/. /theme/

RUN npm ci
RUN npm install -g sass
RUN sass theme.scss:theme.css --style compressed

FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN ls 

RUN CGO_ENABLED=0 GOOS=linux go build -o ddns

FROM gcr.io/distroless/static-debian12 AS run

WORKDIR /

# copy go-build
COPY --from=build --chown=65532:65532 /app/ddns /app/ddns
# copy static files
COPY ./static /app/static
# copy generated css and override png and favicon if available
COPY --from=sass --chown=65532:65532 /theme/theme.css /app/static/css/theme.css
COPY --from=sass --chown=65532:65532 /theme/theme.css.map /app/static/css/theme.css.map
COPY --from=sass --chown=65532:65532 /theme/logo.png* /app/static/img/logo.png
COPY --from=sass --chown=65532:65532 /theme/favicon.ico* /app/static/img/favicon.ico
USER nonroot:nonroot

ENV GIN_MODE=release
ENV DDNS_EXPIRATION_DAYS=90
ENV DDNS_TTL=300

CMD ["/app/ddns"] 
