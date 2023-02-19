package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"listener.hopertz.me/webhooks"
)

func verifyWebhookHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query()

	mode := key.Get("hub.mode")
	token := key.Get("hub.verify_token")
	challenge := key.Get("hub.challenge")

	if len(mode) > 0 && len(token) > 0 {
		if mode == "subscribe" && token == "mytesttoken" {
			w.WriteHeader(http.StatusOK)
			resJson, _ := json.Marshal(challenge)
			w.Write(resJson)
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
	return
}

func webHookEventHandler(w http.ResponseWriter, r *http.Request) {
	var notification webhooks.Notification
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
	router.HandlerFunc(http.MethodGet, "/webhooks", verifyWebhookHandler)
	router.HandlerFunc(http.MethodPost, "/webhooks", webHookEventHandler)

	log.Fatal(http.ListenAndServe(":8080", router))

}
