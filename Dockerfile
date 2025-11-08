FROM golang:1.25 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o /dgo-bot .

FROM scratch
COPY --from=build /dgo-bot /dgo-bot
CMD ["/dgo-bot"]