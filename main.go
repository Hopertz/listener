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

// func HandleNotificationError(ctx context.Context, writer http.ResponseWriter, request *http.Request, err error) error {
// 	if err != nil {
// 		log.Printf("HandleError: %+v\n", err)
// 		return err
// 	}

// 	log.Printf("HandleError: NIL")
// 	return nil
// }



func HandleGeneralNotification(ctx context.Context, writer http.ResponseWriter,notification *hooks.Notification, notificationErrorHandler hooks.NotificationErrorHandler) error {
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
	listenerOptFunc := hooks.WithGenericNotificationHandler(HandleGeneralNotification)
	listener := hooks.NewEventListener()
	listenerOptFunc(listener)
	router.Handler(http.MethodPost, "/webhooks", listener.Handle())

	log.Fatal(http.ListenAndServe(":8080", router))

}
