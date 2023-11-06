package api

import (
	"encoding/json"
	"net/http"
)

func writeEncodedResponse(w http.ResponseWriter, status int, data any) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "json/application")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func decodeRequest(w http.ResponseWriter, r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
