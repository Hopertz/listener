package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	hooks "github.com/piusalfred/whatsapp/webhooks"
)

type verifier struct {
	secret string
	logger io.Writer
	// other places where you pull the secret e.g database
	// other fields for tracing and logging etc

}

// This is first implementation
func (v *verifier) Verify(ctx context.Context, vr *hooks.VerificationRequest) error {
	if vr.Token != v.secret {
		log.Println("invalid token")
		return errors.New("invalid token")
	}

	log.Println("valid token")
	return nil
}

// This is second implementation
func VerifyFn(secret string) hooks.SubscriptionVerifier {
	return func(ctx context.Context, vr *hooks.VerificationRequest) error {
		if vr.Token != secret {
			log.Println("invalid token")
			return errors.New("invalid token")
		}
		log.Println("valid token")
		return nil
	}
}


func main() {
	router := httprouter.New()

	// verifyHandler1
	verifyHandler1 := hooks.VerifySubscriptionHandler(VerifyFn("mytesttoken"))
	router.Handler(http.MethodGet, "/webhooks", verifyHandler1)
	/*
		// verifyHandler2
		verifier := &verifier{
			secret: "mytesttoken",
			logger: log.Writer(),
		}

		verifyHandler2 := hooks.VerifySubscriptionHandler(verifier.Verify)
		router.Handler(http.MethodGet, "/webhooks", verifyHandler2)

	*/
	listener := hooks.NewEventListener().handle()
	router.Handler(http.MethodPost, "/webhooks", listener)

	log.Fatal(http.ListenAndServe(":8080", router))

}
