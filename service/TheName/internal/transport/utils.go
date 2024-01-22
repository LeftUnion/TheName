package transport

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

func writeOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func writeLinkHeader(w http.ResponseWriter, nextURL string) {
	const (
		header = "link"
	)

	w.Header().Set(header, nextURL)
}

func writeJSONResponse(w http.ResponseWriter, obj interface{}, status int) error {
	const (
		header = "content-type"
		value  = "application/json"
	)

	w.Header().Set(header, value)
	w.WriteHeader(status)

	encoder := jsoniter.NewEncoder(w)
	if err := encoder.Encode(obj); err != nil {
		logrus.Errorf("Marshal json: %s", err)
		return err
	}

	return nil
}

func parseJSON(r *http.Request, obj interface{}) error {
	decoder := jsoniter.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()
	decoder.More()
	if err := decoder.Decode(obj); err != nil {
		logrus.Errorf("Декодинг json: %s", err)
		return err
	}

	return nil
}
