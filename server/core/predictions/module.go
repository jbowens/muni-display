package predictions

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jbowens/muni/server/core/config"
	"github.com/octavore/naga/service"
)

const (
	checkInterval = time.Second
)

var (
	// defaultPredicate configures how frequently to query the prediction source for
	// new predictions.
	defaultPredicate = composite(
		atNight(interval(time.Minute)),       // only once per minute at night
		onWeekends(interval(30*time.Second)), // only twice per minute on the weekends
		inMornings(interval(10*time.Second)), // every 10 seconds on weekdays in the morning
		interval(20*time.Second),             // every 20 seconds every other time
	)
)

// Module provides MUNI departure predictions.
type Module struct {
	Config *config.Module

	mu                   sync.Mutex
	keys                 map[string]string
	stops                map[string]Stop
	latestPredictions    map[string][]Prediction
	lastUpdatedTimestamp time.Time
	ticker               *time.Ticker
	predictor            Predictor
}

// Init implements the service.Module interface and installs appropriate lifecycle hooks.
func (m *Module) Init(c *service.Config) {
	c.Setup = m.setup
	c.Start = m.start
}

func (m *Module) setup() error {
	m.keys = make(map[string]string)
	m.stops = make(map[string]Stop)
	m.latestPredictions = make(map[string][]Prediction)

	if err := m.Config.Load("stops.json", &m.stops); err != nil {
		return err
	}
	if err := m.Config.Load("keys.json", &m.keys); err != nil {
		return err
	}

	if _, ok := m.keys["511.org"]; !ok {
		return errors.New("No 511.org access token provided in keys.json")
	}

	m.predictor = &defaultPredictor{accessToken: m.keys["511.org"]}
	return nil
}

func (m *Module) start() {
	fmt.Println("Watching predictions for stops:")
	for k, s := range m.stops {
		fmt.Printf(" - %s (%s %s)\n", k, s.Name, s.Direction)
	}
	if err := m.refreshPredictions(); err != nil {
		panic(err)
	}
	m.ticker = time.NewTicker(checkInterval)
	go m.updatePeriodically()
}

// Current returns all the current route predictions for the given stop.
func (m *Module) Current(stop string) []Prediction {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.latestPredictions[stop]
}

func (m *Module) LastUpdated() time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.lastUpdatedTimestamp
}

func (m *Module) refreshPredictions() error {
	m.mu.Lock()
	m.lastUpdatedTimestamp = time.Now()
	m.mu.Unlock()

	for k, s := range m.stops {
		if err := m.refreshPredictionsForStop(k, &s); err != nil {
			return err
		}
	}
	return nil
}

func (m *Module) refreshPredictionsForStop(key string, stop *Stop) error {
	predictions, err := m.predictor.Predict(stop)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.latestPredictions[key] = predictions
	return nil
}

// updatePeriodically runs in its own goroutine and periodically fetches new departure
// predictions.
func (m *Module) updatePeriodically() {
	for _ = range m.ticker.C {
		shouldUpdate := defaultPredicate(time.Now(), m.LastUpdated(), m)
		if shouldUpdate != nil && *shouldUpdate {
			if err := m.refreshPredictions(); err != nil {
				fmt.Fprintf(os.Stderr, "Error refreshing: %s\n", err.Error())
			}

			// Print the current predictions
			for s := range m.stops {
				var minutes []string
				for _, prediction := range m.Current(s) {
					minutes = append(minutes, strconv.Itoa(prediction.Minutes))
				}
				fmt.Printf("%s: %s\n", s, strings.Join(minutes, ", "))
			}
		}
	}
}
