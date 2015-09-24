package predictions

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	serviceURL = "http://services.my511.org/Transit2.0/GetNextDeparturesByStopCode.aspx?token=&stopCode="
)

type route struct {
	Name          string `xml:"Name,attr"`
	Code          string `xml:"Code,attr"`
	DepartureTime []int  `xml:"RouteDirectionList>RouteDirection>StopList>Stop>DepartureTimeList>DepartureTime"`
}

type nextDeparaturesResponse struct {
	Routes []route `xml:"AgencyList>Agency>RouteList>Route"`
}

type defaultPredictor struct {
	accessToken string
}

var _ Predictor = &defaultPredictor{}

func (d defaultPredictor) Predict(stop *Stop) ([]Prediction, error) {
	l, err := url.Parse(serviceURL)
	if err != nil {
		return nil, err
	}

	queryParams := url.Values{}
	queryParams.Add("token", d.accessToken)
	queryParams.Add("stopCode", strconv.Itoa(int(stop.Code)))
	l.RawQuery = queryParams.Encode()

	resp, err := http.DefaultClient.Post(l.String(), "", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("The 511.org API responded with a non-200 status code: %v", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response nextDeparaturesResponse
	if err := xml.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	var predictions []Prediction
	for _, route := range response.Routes {
		for _, minutes := range route.DepartureTime {
			predictions = append(predictions, Prediction{
				CreatedAt: time.Now(),
				Minutes:   minutes,
				Stop:      stop,
				Source:    "511.org",
			})
		}
	}
	return predictions, nil
}
