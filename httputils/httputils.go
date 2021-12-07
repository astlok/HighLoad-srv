package httputils

import (
	"encoding/json"
	"github.com/prometheus/common/log"
	"highload-srv/metric"
	"net/http"
)

func Respond(w http.ResponseWriter, r *http.Request, requestID uint64, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			RespondError(w, r, requestID, err, http.StatusInternalServerError)
			return
		}
	}
	metric.CrateRequestHits(code, r)
	log.Info(requestID, code)
}

func RespondError(w http.ResponseWriter, r *http.Request, requestID uint64, err error, code int) {
	log.Error(requestID, err)
	metric.CrateRequestError(err)
	Respond(w, r, requestID, code, err)
}
