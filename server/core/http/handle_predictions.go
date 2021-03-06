package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jbowens/muni-display/server/core/predictions"
)

type HandlePredictionsResponse struct {
	LastRefresh time.Time                `json:"last_refresh"`
	Stop        predictions.Stop         `json:"stop"`
	Predictions []predictions.Prediction `json:"predictions"`
}

func (m *Module) handlePredictions(rw http.ResponseWriter, req *http.Request) {
	var zeroStop predictions.Stop

	stopKey := filepath.Base(req.URL.Path)
	stop := m.Predictions.Stop(stopKey)
	if stop == zeroStop {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	lastRefreshed := m.Predictions.LastUpdated()
	predictions := m.Predictions.Current(stopKey)
	m.writeJSON(rw, HandlePredictionsResponse{
		LastRefresh: lastRefreshed,
		Stop:        stop,
		Predictions: predictions,
	})
}

func (m *Module) writeJSON(rw http.ResponseWriter, obj interface{}) {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling: %s\n", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(b)
}
