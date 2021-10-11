package main

import (
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"libs.altipla.consulting/cloudrun"
	"libs.altipla.consulting/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(errors.Stack(err))
	}
}

func run() error {
	var flagDomain, flagTopic string
	flag.StringVarP(&flagDomain, "domain", "d", "", "Domain configured for this webhook to verify the payloads.")
	flag.StringVarP(&flagTopic, "topic", "t", "mailgun-events", "Topic name where the event will be published.")
	flag.Parse()

	if flagDomain == "" {
		return errors.Errorf("--domain flag is required")
	}

	r := cloudrun.Web()
	r.Post("/webhook", WebhookHandler(flagDomain, flagTopic))
	r.Serve()

	return nil
}
