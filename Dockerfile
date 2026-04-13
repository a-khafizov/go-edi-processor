FROM golang:1.25-alpine AS builder
RUN apk add --no-cache make git

WORKDIR /app
COPY . .
RUN go mod download
RUN make build

FROM alpine:latest AS runtime
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/
COPY --from=builder /app/go-doc-history .
COPY --from=builder /app/cmd/.env .
EXPOSE 8080 9090
CMD ["./go-doc-history"]