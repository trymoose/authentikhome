FROM ghcr.io/home-assistant/home-assistant:stable AS builder

COPY --from=golang:1.21-alpine /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o /authentik .

FROM ghcr.io/home-assistant/home-assistant:stable

WORKDIR /authentik
COPY --from=builder /authentik .
RUN touch config.yml