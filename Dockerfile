FROM golang:1.22-alpine

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid 65532 \
  user

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o /app/main .

USER user:user

CMD ["/app/main"]
