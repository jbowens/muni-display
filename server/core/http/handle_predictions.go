package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jbowens/muni/server/core/predictions"
)

type handlePredictionsResponse struct {
	LastRefresh time.Time                `json:"last_refresh"`
	Predictions []predictions.Prediction `json:"predictions"`
}

func (m *Module) handlePredictions(rw http.ResponseWriter, req *http.Request) {
	stopKey := filepath.Base(req.URL.Path)

	lastRefreshed := m.Predictions.LastUpdated()
	predictions := m.Predictions.Current(stopKey)
	if predictions == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	m.writeJSON(rw, handlePredictionsResponse{
		LastRefresh: lastRefreshed,
		Predictions: predictions,
	})
}

func (m *Module) writeJSON(rw http.ResponseWriter, obj interface{}) {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling: %s\n", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(b)
}
