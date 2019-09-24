FROM golang:1.13 AS build

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"'

FROM scratch
ENV PORT 4080
EXPOSE 4080
COPY --from=build /app/pokewants .
ENTRYPOINT ["./pokewants"]
