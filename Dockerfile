FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/blog .

FROM alpine:3.19

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN mkdir -p /app/articles && chown appuser:appgroup /app/articles

COPY --from=builder /out/blog /app/blog
COPY --from=builder /src/templates /app/templates
COPY --from=builder /src/articles /app/articles

ENV PORT=8080
ENV JWT_SECRET=your-secret-key
ENV ADMIN_USERNAME=admin
ENV ADMIN_PASSWORD=your-admin-password

VOLUME ["/app/articles"]

EXPOSE 8080

USER appuser

CMD ["/app/blog"]
