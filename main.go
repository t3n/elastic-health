package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type status struct {
	Status string `json:"status"`
}

func main() {
	zerolog.TimeFieldFormat = ""
	zerolog.LevelFieldName = "severity"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if os.Getenv("ELASTIC_API") == "" {
		log.Fatal().Msg("Set ELASTIC_API environment variable!")
	}
	if os.Getenv("ELASTIC_STATUS") == "" {
		log.Fatal().Msg("Set ELASTIC_STATUS environment variable!")
	}

	http.HandleFunc("/", getStatus)
	port := "8080"
	log.Info().Msgf("Starting Server on Port %s", port)
	log.Fatal().Err(http.ListenAndServe(":"+port, nil))

}

func getStatus(w http.ResponseWriter, r *http.Request) {

	elasticAPI := os.Getenv("ELASTIC_API")
	elasticStatus := os.Getenv("ELASTIC_STATUS")

	elasticHealth := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, elasticAPI, nil)
	if err != nil {
		log.Error().Err(err)
		return
	}

	req.Header.Set("User-Agent", "elastic-healthcheck")

	res, getErr := elasticHealth.Do(req)
	if getErr != nil {
		log.Error().Msg("Elasticsearch is not reachable at" + " " + elasticAPI)
		return
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Error().Err(readErr)
		return
	}

	status := status{}
	jsonErr := json.Unmarshal(body, &status)
	if jsonErr != nil {
		log.Error().Err(jsonErr)
		return
	}

	if status.Status == elasticStatus || status.Status == "green" {
		w.WriteHeader(200)
		return
	}

	w.WriteHeader(503)
}
