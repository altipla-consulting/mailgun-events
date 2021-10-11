package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/mailgun/mailgun-go/v3"
	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"
	"libs.altipla.consulting/pubsub"
	"libs.altipla.consulting/routing"
)

var (
	client *pubsub.Client
)

func init() {
	var err error
	client, err = pubsub.NewClient()
	if err != nil {
		log.Fatal(err)
	}
}

type mailgunEventTags struct {
	Tags []string `json:"tags"`
}

func WebhookHandler(domain, topic string) routing.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var payload mailgun.WebhookPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			return errors.Trace(err)
		}

		log.Infof("%#v", payload)

		mg := mailgun.NewMailgun(domain, os.Getenv("MAILGUN_KEY"))
		signature := mailgun.Signature{
			TimeStamp: payload.Signature.TimeStamp,
			Token:     payload.Signature.Token,
			Signature: payload.Signature.Signature,
		}
		if verified, err := mg.VerifyWebhookSignature(signature); err != nil {
			return errors.Trace(err)
		} else if !verified {
			return errors.Wrapf(err, "failed signature verification: timestamp <%s>, token <%s>, signature <%s>", signature.TimeStamp, signature.Token, signature.Signature)
		}

		tags := new(mailgunEventTags)
		if err := json.Unmarshal(payload.EventData, tags); err != nil {
			return errors.Trace(err)
		}
		log.Infof("%#v", tags)

		attrs := []pubsub.PublishOption{}
		for _, name := range tags.Tags {
			attrs = append(attrs, pubsub.WithAttribute(name, "true"))
		}
		if len(attrs) > 0 {
			if err := client.Topic(topic).PublishBytes(r.Context(), payload.EventData, attrs...); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	}
}
