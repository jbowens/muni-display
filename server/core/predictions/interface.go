package predictions

import "time"

// Predictor defines an interface for things that can predict muni arrival
// times.
type Predictor interface {
	Predict(stop *Stop) ([]Prediction, error)
}

// Prediction encapsulates information about a predicted muni departure from
// a stop.
type Prediction struct {
	CreatedAt time.Time
	Minutes   int
	Stop      *Stop
	Source    string
}

// Stop represents a public-transit stop and the information required to query
// prediction data for the stop.
type Stop struct {
	Agency    string `json:"agency"`
	Route     string `json:"route"`
	Direction string `json:"direction"`
	Name      string `json:"name"`
	Code      int    `json:"code"`
}
