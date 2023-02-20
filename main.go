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
	"listener.hopertz.me/webhooks"
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

type handler struct {
}

func (h *handler) HandleError(ctx context.Context, writer http.ResponseWriter, request *http.Request, err error) error {
	if err != nil {
		log.Printf("HandleError: %+v\n", err)
		return err
	}

	log.Printf("HandleError: NIL")
	return nil
}

func (h *handler) HandleEvent(ctx context.Context, writer http.ResponseWriter, request *http.Request, notification *hooks.Notification) error {
	os.Stdout.WriteString("HandleEvent")
	jsonb, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	// print the string representation of the json
	//os.Stdout.WriteString(string(jsonb))
	log.Printf("\n%s\n", string(jsonb))
	writer.WriteHeader(http.StatusOK)
	return nil
}

func verifyWebhookHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query()

	mode := key.Get("hub.mode")
	token := key.Get("hub.verify_token")
	challenge := key.Get("hub.challenge")

	if len(mode) > 0 && len(token) > 0 {
		if mode == "subscribe" && token == "mytesttoken" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(challenge))
			log.Println("Hurray")
			return

		} else {
			w.WriteHeader(http.StatusForbidden)
			log.Println("forbiden")
			return
		}

	}
	w.WriteHeader(http.StatusBadRequest)
	log.Println("bad request")
}

func webHookEventHandler(w http.ResponseWriter, r *http.Request) {
	var notification webhooks.WebhookMessage
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&notification)
	log.Println(notification)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	} else {
		w.WriteHeader(http.StatusAccepted)
	}

}

func main() {
	router := httprouter.New()
	/*
		// verifyHandler1
		verifyHandler1 := hooks.VerifySubscriptionHandler(VerifyFn("mytesttoken"))
		router.Handler(http.MethodGet, "/webhooks", verifyHandler1)
		// verifyHandler2
		verifier := &verifier{
			secret: "mytesttoken",
			logger: log.Writer(),
		}

		verifyHandler2 := hooks.VerifySubscriptionHandler(verifier.Verify)
		router.Handler(http.MethodGet, "/webhooks", verifyHandler2)
	*/

	router.HandlerFunc(http.MethodGet, "/webhooks", verifyWebhookHandler)

	// This is the Event Handler Implementation
	// handler := &handler{}
	// listener := hooks.NewEventListener(handler)
	// router.Handler(http.MethodPost, "/webhooks", listener)

	router.HandlerFunc(http.MethodPost, "/webhooks", webHookEventHandler)

	log.Fatal(http.ListenAndServe(":8080", router))

}
