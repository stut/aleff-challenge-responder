package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"net/http"
	"os"
)

var UrlPrefix string
var KvRoot string

func sendChallengeResponse(w http.ResponseWriter, req *http.Request) {
	// This is safe because we wouldn't be in this function if the URL doesn't have this prefix.
	token := req.URL.Path[len(UrlPrefix):]

	if token == "health" {
		w.WriteHeader(204)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		w.WriteHeader(500)
		_, _ = fmt.Fprintf(w, "Backend connection failed!")
		log.Printf("Failed to connect to consul: %v", err)
		return
	}

	kv := client.KV()

	var val *api.KVPair
	val, _, err = kv.Get(fmt.Sprintf("%s%s", KvRoot, token), nil)
	if err != nil {
		w.WriteHeader(500)
		_, _ = fmt.Fprintf(w, "Failed to fetch token!")
		log.Printf("Failed to fetch token [%s]: %v", token, err)
		return
	}

	if val == nil {
		w.WriteHeader(404)
		_, _ = fmt.Fprintf(w, "Token not found!")
		log.Printf("Token not found: %s", token)
		return
	}

	_, _ = fmt.Fprint(w, string(val.Value))
}

func main() {
	listenPort := os.Getenv("NOMAD_PORT_http")
	if len(listenPort) == 0 {
		listenPort = "8080"
	}

	UrlPrefix = os.Getenv("URL_PREFIX")
	if len(UrlPrefix) == 0 {
		UrlPrefix = "/.well-known/acme-challenge/"
	}

	KvRoot = os.Getenv("KV_ROOT")
	if len(KvRoot) == 0 {
		KvRoot = "certs/challenges/"
	}

	http.HandleFunc(UrlPrefix, sendChallengeResponse)
	panic(http.ListenAndServe(fmt.Sprintf(":%s", listenPort), nil))
}
