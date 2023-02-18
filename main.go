package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"listener.hopertz.me/webhooks"
)

type Middleware func(next http.Handler) http.Handler

// Wraps a http.Handler with a middlewares
func Wrap(handler http.Handler, middlewares ...Middleware) http.Handler {
	// wraps backwards
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// StartTimeMiddleware is a middleware that adds the start time to the request context
func StartTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("StartTimeMiddleware")
		ctx := r.Context()
		ctx = context.WithValue(ctx, "startTime", time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// EndTimeMiddleware is a middleware that adds the end time to the request context
func EndTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("EndTimeMiddleware")
		ctx := r.Context()
		ctx = context.WithValue(ctx, "endTime", time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

var _ webhooks.EventHandler = (*handler)(nil)

type handler struct {
	// it can be db conections, etc
	// connection to notify the user
	// etc etc
}

func (h *handler) HandleError(ctx context.Context, writer http.ResponseWriter, request *http.Request, err error) error {
	if err != nil {
		os.Stdout.WriteString(err.Error())
	}

	os.Stdout.WriteString("error is nil")
	return nil
}

func (h *handler) HandleEvent(ctx context.Context, writer http.ResponseWriter, request *http.Request, notification *webhooks.Notification) error {
	os.Stdout.WriteString("HandleEvent")
	jsonb, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(jsonb)
	return nil
}

func main() {
	// Create a new handler
	handler := &handler{}
	ls := webhooks.NewEventListener(handler)
	middlewares := []Middleware{
		StartTimeMiddleware,
		EndTimeMiddleware,
	}
	finalHandler := Wrap(ls, middlewares...)
	mux := http.NewServeMux()
	mux.Handle("/webhooks", finalHandler)

	// Create a new server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

}
