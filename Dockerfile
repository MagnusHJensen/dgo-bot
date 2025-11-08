FROM golang:1.25 AS build

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o /dgo-bot .

FROM gcr.io/distroless/static:nonroot
COPY --from=build /dgo-bot /dgo-bot

# Run as nonroot (safe default user included in distroless)
USER nonroot:nonroot

ENTRYPOINT ["/dgo-bot"]