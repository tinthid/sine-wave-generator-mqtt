FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o demo_iot_sensor

FROM gcr.io/distroless/static

COPY --from=builder /app/demo_iot_sensor /

CMD ["/demo_iot_sensor"]
