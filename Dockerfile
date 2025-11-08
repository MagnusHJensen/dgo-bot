FROM golang:1.25 as build
WORKDIR /app
COPY . .
RUN go build -o /dgo-bot .

FROM scratch
COPY --from=build /dgo-bot /dgo-bot
CMD ["/dgo-bot"]