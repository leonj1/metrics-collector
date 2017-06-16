package routes

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"metrics-collector/models"
	"net/http"
	"strconv"
	"time"
)

func GetMetrics(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var m models.MetricValue

	name := ps.ByName("name")
	if name == "" {
		http.Error(w, "No name provided", http.StatusBadRequest)
		return
	}

	days := ps.ByName("days")
	if days == "" {
		http.Error(w, "Days not provided", http.StatusBadRequest)
		return
	}

	numDays, err := strconv.ParseInt(days, 10, 64)
	if err != nil {
		http.Error(w, "Invalid days", http.StatusBadRequest)
		return
	}
	numDays = numDays * -1

	metrics, err := m.FindByMetricNameBetweenDates(
		name,
		time.Now().AddDate(0,0, int(numDays)),
		time.Now(),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(&metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(ContentType, JSON)
	w.Write(js)
}