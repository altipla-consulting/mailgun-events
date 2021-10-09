
FROM golang:1.11 as builder

WORKDIR /opt/ac
COPY . .

RUN go install ./cmd/mailgun-events

# ==============================================================================

FROM gcr.io/distroless/base

COPY --from=builder /build/mailgun-events .

ENTRYPOINT [ "./mailgun-events" ]
