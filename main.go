package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/kenethrrizzo/banking-auth/app"
)

func main() {
	log.Info("Starting application")
	app.Start()
}
