package main

import (
	"io"
	"net/http"
)

/*
 * Modify this constant to one of the possible enum values
 * and observe whether the USB mounts!
 */
const currentCompetitionStatus = CompetitionOngoing

// CompetitionStatus determines the current status of the competition.
type CompetitionStatus string

const (
	CompetitionNotStarted CompetitionStatus = "not_started"
	CompetitionOngoing    CompetitionStatus = "ongoing"
	CompetitionFinished   CompetitionStatus = "finished"
)

func main() {
	m := http.NewServeMux()
	m.HandleFunc("/competition-status", competitionStatusHandler)

	println("Competition status is " + string(currentCompetitionStatus))
	println("Listening to localhost:8080")
	if err := http.ListenAndServe(":8080", m); err != nil {
		panic(err)
	}
}

func competitionStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, string(currentCompetitionStatus))
}

