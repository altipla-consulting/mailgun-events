
FROM golang:1.17 as builder

WORKDIR /opt/ac

COPY cmd cmd
COPY go.mod go.mod
COPY go.sum go.sum

RUN go install ./cmd/mailgun-events

# ==============================================================================

FROM gcr.io/distroless/base

COPY --from=builder /go/bin/mailgun-events .

ENTRYPOINT [ "./mailgun-events" ]
