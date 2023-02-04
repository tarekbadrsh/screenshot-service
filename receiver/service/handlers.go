package service

import (
	"fmt"
	"receiver/config"
	"receiver/logger"
	"receiver/messaging"
	"receiver/parsing"

	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// start initialize topics
var c = config.Configuration()
var rawURLTopic = c.RawURLTopic

// end initialize topics

type responseMsg struct {
	Msg string `json:"message"`
}

func receiveJSON(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h := parsing.GetHandler(parsing.JSON, rawURLTopic)
	err := h.Handler(r.Body)
	if err != nil {
		logger.Error(err)
		msg := fmt.Sprintf(`{"message":"%v"}`, err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	logger.Info("Produced Successfully")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseMsg{"Success"})
}

func ready(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if messaging.KafkaReadiness {
		json.NewEncoder(w).Encode(responseMsg{"Ready"})
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, `{"message":"Unready"}`, http.StatusServiceUnavailable)
}

// logmid : logging midleware
func logmid(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		logger.Infof("[%s] on: %s", r.Method, r.URL)
		next(w, r, ps)
	}
}
